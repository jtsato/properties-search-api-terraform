@echo on

terraform plan -destroy

terraform destroy -auto-approve -var-file=secrets.tfvars

:: Just a pause to see if everything is destroyed on the Google Cloud Console
pause

terraform state list

terraform plan -out=tfplan.out -var-file=secrets.tfvars

terraform apply -auto-approve tfplan.out
