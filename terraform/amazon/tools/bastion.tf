variable "bastion_instance_type" {
  default = "t2.nano"
}

variable "bastion_ami" {}

variable "bastion_root_size" {
  default = "16"
}

resource "aws_eip" "bastion" {
  vpc      = true
  instance = "${aws_instance.bastion.id}"
}

resource "aws_security_group" "bastion" {
  name        = "${data.template_file.stack_name.rendered}-bastion"
  vpc_id      = "${data.terraform_remote_state.network.vpc_id}"
  description = "Bastion instance in ${data.template_file.stack_name.rendered}"

  tags {
    Name        = "${data.template_file.stack_name.rendered}-bastion"
    Environment = "${var.environment}"
    Project     = "${var.project}"
    Contact     = "${var.contact}"
  }
}

resource "aws_security_group_rule" "egress_allow_all" {
  type              = "egress"
  protocol          = -1
  from_port         = 0
  to_port           = 65535
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = "${aws_security_group.bastion.id}"
}

resource "aws_security_group_rule" "ingress_allow_ssh_all" {
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 22
  to_port           = 22
  cidr_blocks       = ["${var.admin_ips}"]
  security_group_id = "${aws_security_group.bastion.id}"
}

data "template_file" "bastion_user_data" {
  template = "${file("${path.module}/templates/bastion_user_data.yaml")}"

  vars {
    fqdn = "bastion.${data.terraform_remote_state.state.public_zone}"
  }
}

resource "aws_instance" "bastion" {
  ami                    = "${var.bastion_ami}"
  instance_type          = "${var.bastion_instance_type}"
  subnet_id              = "${data.terraform_remote_state.network.public_subnet_ids[0]}"
  key_name               = "${var.key_name}"
  vpc_security_group_ids = ["${aws_security_group.bastion.id}"]

  root_block_device = {
    volume_type = "gp2"
    volume_size = "${var.bastion_root_size}"
  }

  tags {
    Name        = "${data.template_file.stack_name.rendered}-bastion"
    Environment = "${var.environment}"
    Project     = "${var.project}"
    Contact     = "${var.contact}"
    tarmak_role = "bastion"
  }

  user_data = "${data.template_file.bastion_user_data.rendered}"
}

resource "tarmak_bastion_instance" {
  username = "centos"
  ip = "${aws_eip.bastion.public_ip}"
  depend_on = ["${aws_instance.bastion}"]
}

resource "aws_route53_record" "bastion" {
  zone_id = "${data.terraform_remote_state.state.public_zone_id}"
  name    = "bastion.${var.environment}"
  type    = "A"
  ttl     = "300"
  records = ["${aws_eip.bastion.public_ip}"]
}

resource "aws_route53_record" "bastion_private" {
  zone_id = "${data.terraform_remote_state.network.private_zone_id.0}"
  name    = "bastion.${var.environment}"
  type    = "A"
  ttl     = "60"
  records = ["${aws_instance.bastion.private_ip}"]
}

output "bastion_fqdn" {
  value = "${aws_route53_record.bastion.fqdn}"
}

output "bastion_private_ip" {
  value = "${aws_eip.bastion.public_ip}"
}

output "bastion_ip" {
  value = "${aws_eip.bastion.public_ip}"
}

output "bastion_security_group_id" {
  value = "${aws_security_group.bastion.id}"
}
