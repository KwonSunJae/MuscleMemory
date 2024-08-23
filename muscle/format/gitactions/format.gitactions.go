package format

const (
	GitActionsAnsibleTemplateWithPW = `
on:
  push:
    branches:
      - {{.project-name}}

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
      rm -rf {{.project-name}}
      git clone -b {{.project-name}} {{.repository}} {{.project-name}}
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
    - {{.project-name}}
  
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
    rm -rf {{.project-name}}
    git clone -b {{.project-name}} {{.repository}} {{.project-name}}
    cd {{.project-name}}
    for playbook in $(find . -name "*.yml"); do
           ansible-playbook -i inventory "$playbook"
    done
    cd ..
`
	GitAcitonsTerraformTemplateWithPW = `
on:
  push:
    branches:
    - {{.project-name}}
  
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
    rm -rf {{.project-name}}
    git clone -b {{.project-name}} {{.repository}} {{.project-name}}
    cd {{.project-name}}
    terraform init
    terraform apply -auto-approve
    cd ..
`
	GitAcitonsTerraformTemplateWithKey = `
on:
  push:
    branches:
    - {{.project-name}}

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
    rm -rf {{.project-name}}
    git clone -b {{.project-name}} {{.repository}} {{.project-name}}
    cd {{.project-name}}
    terraform init
    terraform apply -auto-approve
    cd ..
`
	GitActionsTerraformPublicTemplate = `
name: Terraform Apply

on:
  push:
    branches:
      - {{ .branch }}

jobs:
  terraform:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout code
        uses: actions/checkout@v2
      - name: Set up Terraform
        uses: HashiCorp/setup-terraform@v1
        with:
          terraform_version: 1.5.7

      - name: Terraform Init
        run: terraform init

      - name: Terraform Plan
        id: plan
        run: terraform plan -out=tfplan

      - name: Terraform Apply
        run: terraform apply -auto-approve tfplan

      - name: Commit tfstate file
        run: |
          rm -rf .lock
          git config --local user.email "{{ .useremail }}"
          git config --local user.name "{{ .username }}"
          git add .
          git commit -m "Update tfstate file"
          git push
        env:
          GITHUB_TOKEN: {{ .github_token}}
`
	GitActionsAnsiblePublicTemplate = `
name: Ansible Playbook Runner

on:
  push:
    branches:
      - {{.project-name}}

jobs:
  ansible:
    runs-on: ubuntu-latest

steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Set up Python
      uses: actions/setup-python@v2
      with:
      python-version: '3.8'  # 사용할 Python 버전

    - name: Install Ansible
      run: |
        python -m pip install --upgrade pip
        pip install ansible

    - name: Run Ansible Playbook
      run: ansible-playbook {{.playbookname}} 
      env:
        ANSIBLE_HOST_KEY_CHECKING: "False"
`
)
