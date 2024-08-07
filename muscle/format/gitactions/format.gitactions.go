package format

const (
	GitActionsAnsibleTemplateWithPW = `
on:
  push:
    branches:
      - {{.branch}}

jobs:
  deploy:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout repository
      uses: actions/checkout@v3

    - name: Prepare environment
      uses: appleboy/ssh-action@master
      with:
        host: {{.host}}
        username: {{.username}}
        password: {{.password}}
        script: |
		  echo "{{.rootpassword}}" | sudo -S apt-get update -y && sudo apt-get upgrade -y
		  sudo apt-get install -y ansible git
		  rm -rf {{.branch}}
		  git clone -b {{.branch}} {{.repository}} {{.branch}}
		  cd {{.branch}}
		  for playbook in $(find . -name "*.yml"); do
   		  	ansible-playbook -i inventory "$playbook"
		  done
		  cd ..
`
	GitActionsAnsibleTemplateWithKey = `
on:
  push:
    branches:
	  - {{.branch}}
	
jobs:
  deploy:
    runs-on: ubuntu-latest

	steps:
	- name: Checkout repository
	  uses: actions/checkout@v3
	  with:
		host: {{.host}}
		username: {{.username}}
		key: {{.key}}
		script: |
		echo "{{.rootpassword}}" | sudo -S apt-get update -y && sudo apt-get upgrade -y
		sudo apt-get install -y ansible git
		rm -rf {{.branch}}
		git clone -b {{.branch}} {{.repository}} {{.branch}}
		cd {{.branch}}
		for playbook in $(find . -name "*.yml"); do
   		  	ansible-playbook -i inventory "$playbook"
		done
		cd ..
`
	GitAcitonsTerraformTemplateWithPW = `
on:
  push:
    branches:
	  - {{.branch}}
	
jobs:
  deploy:
    runs-on: ubuntu-latest

	steps:
	- name: Checkout repository
	  uses: actions/checkout@v3
	  with:
		host: {{.host}}
		username: {{.username}}
		password: {{.password}}
		script: |
		echo "{{.rootpassword}}" | sudo -S apt-get update -y && sudo apt-get upgrade -y
		sudo apt-get install -y terraform git
		rm -rf {{.branch}}
		git clone -b {{.branch}} {{.repository}} {{.branch}}
		cd {{.branch}}
		terraform init
		terraform apply -auto-approve
		cd ..
`
	GitAcitonsTerraformTemplateWithKey = `
on:
  push:
    branches:
	  - {{.branch}}

jobs:
  deploy:
    runs-on: ubuntu-latest

	steps:
	- name: Checkout repository
	  uses: actions/checkout@v3
	  with:
		host: {{.host}}
		username: {{.username}}
		key: {{.key}}
		script: |
		echo "{{.rootpassword}}" | sudo -S apt-get update -y && sudo apt-get upgrade -y
		sudo apt-get install -y terraform git
		rm -rf {{.branch}}
		git clone -b {{.branch}} {{.repository}} {{.branch}}
		cd {{.branch}}
		terraform init
		terraform apply -auto-approve
		cd ..
`
)
