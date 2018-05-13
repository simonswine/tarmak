variable "region" {
}

variable "cluster_name" {
}

provider "aws" {
  region              = "${var.region}"
}

data "aws_caller_identity" "current" {
  provider = "aws"
}

data "template_file" "assume_role_policy" {
   template = <EOS
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    },
    {
      "Action": "sts:AssumeRole",
      "Principal": {
         "AWS": "${kubernetes_worker_role_arn}"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOS

  vars {
    kubernetes_worker_role_arn = "arn:aws:iam::${data.aws_caller_identity.current.account_id}:role/kubernetes.${data.template_file.stack_name.rendered}.kubernetes_worker"
  }
}
