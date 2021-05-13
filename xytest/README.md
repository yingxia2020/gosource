# Install Kubernetes cluster with KOPS on AWS and run CNB


### References:
- https://medium.com/containermind/how-to-create-a-kubernetes-cluster-on-aws-in-few-minutes-89dda10354f4
- https://github.com/kubernetes/kops/blob/master/docs/aws.md
- https://medium.com/@mcyasar/amazon-aws-kubernetes-kops-installation-7a205fe2d118


## Preparations:

#### Install kops
```
curl -LO https://github.com/kubernetes/kops/releases/download/$(curl -s https://api.github.com/repos/kubernetes/kops/releases/latest | grep tag_name | cut -d '"' -f 4)/kops-linux-amd64
chmod +x kops-linux-amd64
sudo mv kops-linux-amd64 /usr/local/bin/kops
kops version
```

#### Install kubectl
```
sudo apt-get update && sudo apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl
```

#### Create AWS IAM user with the following permissions
    - AmazonEC2FullAccess
    - AmazonRoute53FullAccess
    - AmazonS3FullAccess
    - AmazonVPCFullAccess


#### Create access key for the user, then add the key to ~/.aws/credentials file
```
[default]
aws_access_key_id = AKIXXXXXXXXXXXXX
aws_secret_access_key = sIrkzNOXxXXXXXXXXXXXxxxXXXX
```

#### Install AWS CLI, the recommended way by AWS is using pip or pip3
```
sudo apt-get install python3
sudo apt install python3-pip
sudo pip3 install awscli --upgrade --user
aws configure
AWS Access Key ID [None]: AKIXXXXXXXXXXXXX
AWS Secret Access Key [None]: sIrkzNOXxXXXXXXXXXXXxxxXXXX
Default region name [None]: us-west-2
Default output format [None]: json
```

#### Update .bashrc file with the following lines
```
export bucket_name=cnb-kops-state-store
export KOPS_CLUSTER_NAME=cnb.k8s.local
export KOPS_STATE_STORE=s3://${bucket_name}
```

#### Source .bashrc
```
source .bashrc
```

#### Create S3 bucket for KOPS to run
```
aws s3 mb s3://${bucket_name} --region us-west-2
aws s3api put-bucket-versioning --bucket ${bucket_name} --versioning-configuration Status=Enabled
```

#### Make sure you have one spare VPC available for KOPS to run
KOPS will automatically create a VPC for your cluster to run within. If you don't have a spare one available, you will not be able to create your cluster.

#### Create a key pair on AWS console, i.e. ('cnb_aws_key'). Download the save the cnb_aws_key.pem file
**NOTE:** Keep the file in a safe place. Do not lose this file!

#### Create a public key using AWS linux machines’ private key
```
cd ~/.ssh
mkdir test
cp cnb_aws_key.pem ~/.ssh/test
cd ~/.ssh/test
ssh-keygen -y
Enter file in which the key is:  /home/user/.ssh/test/cnb_aws_key.pem
put the content to file: ~/.ssh/test/cnb_aws.pub
```

## Create Cluster

#### Example cluster configuration
```
kops create cluster --node-count=2 --node-size=m5a.24xlarge --master-size=m5a.4xlarge --zones=us-west-2a --name=${KOPS_CLUSTER_NAME}
```
**NOTE:** You can edit these options to create the cluster you'd like

#### Delete the default used by kops for this cluster 
```
kops describe secret
kops delete secret sshpublickey admin   78:18:6b:07:e5:40:14:cb:54:c0:6e:d7:b1:90:9d:43
kops describe secret
```

#### Register the public key created previously with KOPS
```
kops create secret sshpublickey admin -i ~/.ssh/test/cnb_aws.pub --name=${KOPS_CLUSTER_NAME} --state s3://${bucket_name}
```

#### Optional - If need enable metrics server and run workload in HPA mode
Current there is a bug: https://github.com/kubernetes/kops/pull/6201. The current work around:
```
kops edit cluster
kubelet:
    anonymousAuth: false
    authenticationTokenWebhook: true     <--- Add this line
```

#### Deploy cluster
```
kops update cluster --name ${KOPS_CLUSTER_NAME} --yes
```

#### Wait for some time (around 5-10 minutes) and validate cluster
```
kops validate cluster
kubectl get nodes
kubectl cluster-info
```

## Run CNB and Save Results

#### Access cluster

- Go to AWS console
- Go to EC2 running instances page
- Choose instance name "master-us-west-2a.masters.cnb.k8s.local"
- Click "Connect"
- Get the connection string under the Public DNS section, it will have the following format:

    **ec2-54-189-181-18.us-west-2.compute.amazonaws.com**

#### Copy CNB release package to master node
```
scp -i "~/.ssh/test/cnb_aws_key.pem" cnbrun.tar.gz admin@ec2-54-189-181-18.us-west-2.compute.amazonaws.com:~/
```
**Note:** Make sure you use your own connection string!

#### SSH into master node
```
ssh -i "~/.ssh/test/cnb_aws_key.pem" admin@ec2-54-189-181-18.us-west-2.compute.amazonaws.com
```

#### Run CNB
```
gunzip cnbrun.tar.gz
tar xvf cnbrun.tar
cd cnbrun
```

Modify config.json file according to README in cnbrun directory

Run the benchmark:
```
./cnbrun
```

**NOTE:** Results will be written in the 'output' directory.

#### Save results locally

```
scp -i "~/.ssh/test/cnb_aws_key.pem" admin@ec2-54-189-181-18.us-west-2.compute.amazonaws.com:~/cnbrun/output/* .
```

## Clean up Cluster
After you are done running CNB and have saved the results:
```
kops delete  cluster --name ${KOPS_CLUSTER_NAME} --yes
```
