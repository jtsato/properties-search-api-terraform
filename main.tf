resource "null_resource" "app_server" {

  connection {
    type     = "ssh"
    user     = var.ssh_user
    password = var.ssh_password
    host     = var.ssh_host
    port     = var.ssh_port
  }

  provisioner "file" {
    source      = "install-docker.sh"
    destination = "/tmp/install-docker.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "echo ${var.ssh_password} | sudo -S chmod +x /tmp/install-docker.sh",
      "echo ${var.ssh_password} | sudo -S /tmp/install-docker.sh",
    ]
  }
}

# gsutil mb -p duckhome-firebase -c STANDARD -l southamerica-east1 gs://duckhome-pps-terraform-state
terraform {
  backend "gcs" {
    bucket = "duckhome-pps-terraform-state"
    prefix = "terraform/state"
  }
}
