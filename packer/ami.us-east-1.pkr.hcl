packer {
  required_plugins {
    amazon = {
      source  = "github.com/hashicorp/amazon"
      version = ">= 1.0.0"
    }
  }
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "source_ami" {
  type    = string
  default = "ami-06db4d78cb1d3bbf9" # Debian 12
}

variable "ssh_username" {
  type    = string
  default = "admin"
}

variable "subnet_id" {
  type    = string
  default = "subnet-0f61ed28cc7a4aab6"
}

variable "vpc_id" {
  type    = string
  default = "vpc-00a4f3cbc81b90084"
}

variable "demo_user_id" {
  type    = string
  default = "118157167744"
}

variable "app_binary_path" {
  type    = string
  default = "myapp"
}

variable "users_path" {
  type    = string
  default = "users.csv"
}

# https://www.packer.io/plugins/builders/amazon/ebs
source "amazon-ebs" "my-ami" {
  region          = "${var.aws_region}"
  ami_name        = "csye6225_${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  ami_description = "AMI for CSYE 6225"
  ami_users       = ["${var.demo_user_id}"]
  ami_regions = [
    "us-east-1",
  ]

  aws_polling {
    delay_seconds = 120
    max_attempts  = 50
  }

  instance_type = "t2.micro"
  source_ami    = "${var.source_ami}"
  ssh_username  = "${var.ssh_username}"
  subnet_id     = "${var.subnet_id}"
  vpc_id        = "${var.vpc_id}"

  launch_block_device_mappings {
    delete_on_termination = true
    device_name           = "/dev/xvda"
    volume_size           = 8
    volume_type           = "gp2"
  }
}

build {
  sources = ["source.amazon-ebs.my-ami"]

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y golang", # 安装 Golang
      # 创建用户和组
      "sudo groupadd csye6225",
      "sudo useradd -s /bin/false -g csye6225 -d /opt/csye6225 -m csye6225",
    ]
  }

  provisioner "file" {
    source      = "${var.app_binary_path}" # Go 二进制文件的路径
    destination = "/tmp/myapp"             # 临时路径
  }

  provisioner "file" {
    source      = "${var.users_path}"
    destination = "/tmp/users.csv" # 临时路径
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/myapp /opt/csye6225/myapp", # 将文件移动到最终位置
      "sudo mv /tmp/users.csv /opt/csye6225/users.csv",
      "sudo chown csye6225:csye6225 /opt/csye6225/myapp",
      "sudo chown csye6225:csye6225 /opt/csye6225/users.csv",
      "sudo chmod 755 /opt/csye6225/myapp", #rwx r-x r-x
    ]
  }

  provisioner "shell" {
    script = "packer/setup.sh"
  }

  # 将服务文件复制到 systemd 的目录并启动服务
  provisioner "shell" {
    inline = [
      "sudo cp /tmp/csye6225.service /etc/systemd/system",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable csye6225",
      "sudo systemctl start csye6225",
      "sudo chown csye6225:csye6225 /etc/systemd/system/csye6225.service"
    ]
  }
  # 下载cloud watch
  provisioner "shell" {
    inline = [
      "curl -sL https://s3.amazonaws.com/amazoncloudwatch-agent/debian/amd64/latest/amazon-cloudwatch-agent.deb -o amazon-cloudwatch-agent.deb",
      "sudo dpkg -i -E ./amazon-cloudwatch-agent.deb",
      "sudo systemctl enable amazon-cloudwatch-agent",
      "sudo systemctl start amazon-cloudwatch-agent",
    ]
  }

  provisioner "file" {
    source      = "packer/cloudwatch-config.json"
    destination = "/tmp/cloudwatch-config.json"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/cloudwatch-config.json /opt/cloudwatch-config.json",
      "sudo mkdir /var/log/webapp",
      "sudo chown csye6225:csye6225 /var/log/webapp",
    ]
  }

}
