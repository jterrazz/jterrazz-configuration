# VPS Setup Guide

## 1. Hardware Setup (VPS Provider)

1. Order Ubuntu Latest VPS from provider
2. Generate SSH key:
   ```bash
   ssh-keygen -t ed25519 -C "email@example.com"
   ```
3. Save SSH key in Bitwarden
4. Add SSH public key to VPS provider
5. Connect to VPS:
   ```bash
   ssh -i id_hetzner root@VPS_IP
   ```

## 2. Software Setup (Coolify)

1. Access Coolify interface at `VPS_IP:8000`
2. Initial configuration:
   - Setup servers
   - Configure DNS with Cloudflare proxy

## 3. Maintenance

### Credentials (in Bitwarden)

- Copy `/data/coolify/source/.env`
- Save SSH key

### Backups

- Regular backups of Coolify data
