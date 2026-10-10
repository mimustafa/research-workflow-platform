# Ansible VM deployment

The playbook prepares an Ubuntu VM with Docker Engine and Docker Compose, clones this
repository, creates a protected `.env`, and starts the Compose deployment.

## Use

```bash
cd deploy/ansible
cp inventory.ini.example inventory.ini
cp group_vars/all.yml.example group_vars/all.yml
ansible-playbook -i inventory.ini deploy.yml
```

Do not commit `inventory.ini` or `group_vars/all.yml`; they contain host addresses and
deployment credentials. For a real environment, store the credential values in Ansible
Vault or retrieve them from OpenBao.
