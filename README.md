# ✅ Feature List for First Release (MVP)

## 🔐 1. Authentication & Authorization

Centralized auth for all user types.

- User registration & login (email + password)
- Forgot/reset password
- JWT-based session management
- Role-based access control — **Doctor** and **Admin**
- Default role on registration is **Doctor**
- Admin can only be assigned manually (no self-registration)
- Secure password hashing

✅ **Goal:** Secure, role-aware access.

## 👨‍⚕️ 2. Doctor Profile Management

Each doctor manages their own workspace.

- Doctor profile setup after registration
    - Name
    - Degree
    - Specialization
    - Chamber/hospital info
    - Signature upload
- Profile tied to authenticated user
- Doctor can access **only their own** patients & prescriptions

✅ **Goal:** Personalized and secure workspace.


## 👤 3. Patient Management

Core foundation of the system.

- Add new patient
- Edit patient info
- Patient list with search
- View patient profile

**Basic patient fields:**

- Name
- Age / Date of birth
- Gender
- Phone
- Address

✅ **Goal:** Quick patient management.


## 📋 4. Prescription Creation (Core Feature ⭐)

Heart of the product.

- Create prescription for a patient
- Select prescription template
- Add visit date
- Chief complaints
- Diagnosis notes
- Add medicines:
    - Medicine name
    - Dosage
    - Frequency
    - Duration
    - Instructions
- Advice/recommendations
- Next visit date
- Save prescription (immutable after finalization — recommended)

✅ **Goal:** Generate complete digital prescriptions.


## 🧾 5. Prescription Templates

Doctors can customize how prescriptions look.

### From Settings, doctor can:

- View multiple built-in templates
- Preview template
- Set default template
- (Optional MVP+) upload clinic logo

### System provides:

- Multiple professional layouts (2–3 templates for MVP)
- Header with doctor info
- Signature placement
- Footer support

✅ **Goal:** Flexible but simple prescription styling.


## 🖨️ 6. Prescription View, Print & PDF

Doctors MUST be able to print.

- Auto-formatted prescription layout
- Apply selected template
- Doctor header & signature
- Print-friendly view
- Download as PDF
- View past prescriptions of a patient

✅ **Goal:** Professional ready-to-print output.


## 📚 7. Patient History (High Value)

Long-term retention feature.

- Timeline of patient visits
- View previous prescriptions
- View past diagnoses
- Search within patient history

✅ **Goal:** Quick medical history access.


## 🔍 8. Basic Search & Filters

For usability.

- Search patients by name/phone
- Filter prescriptions by date
- Recent patients list

✅ **Goal:** Fast navigation.


## ⚙️ 9. Basic Settings

Doctor-level preferences.

### General settings:

- Chamber timing
- Default prescription footer
- Time zone & date format

### 🧾 Prescription template settings:

- Choose default template
- Preview templates
- Upload clinic logo (optional but recommended)

✅ **Goal:** Doctor-controlled personalization.









### Dev & Tooling
- [golangci-lint](https://github.com/golangci/golangci-lint): Linter aggregator
- [Docker](https://www.docker.com/), [docker-compose](https://docs.docker.com/compose/): Containerization

### Enable git pre-commit hook

- Run the following command. we are using golangci-lint as a linter.

```Shell
make pre-commit-hook 
```

- Perfect! now for each commit, golangci-lint will check for any issue.

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations, integrated with the Cobra CLI in the `doc-mate` service. Migrations are written as simple `.up.sql` and `.down.sql` files.

### Migration Directory
- Migration files are placed in the following directory:
    - `./migrations`

### Running Migrations

You can run migrations using the Cobra CLI (from the built binary):

```sh
# Migrate up (apply all new migrations)
./app migrate up

# Migrate down (revert all migrations)
./app migrate down
```

- **DB credentials/config** are loaded from Consul at runtime.

### Running Migrations Locally (with go run)

You can run migrations directly using `go run` for local development. Make sure to set the required Consul environment variables:

```sh
export CONSUL_URL=http://localhost:8500
export CONSUL_PATH=docmate

go run main.go migrate up
# or to revert all migrations
go run main.go migrate down
```

You can also set the environment variables inline:

```sh
CONSUL_URL=http://localhost:8500 CONSUL_PATH=docmate go run main.go migrate up
```
- DB credentials/config are loaded from Consul at runtime.

### Creating a New Migration

You can create a new migration using the [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate):

1. **Install the CLI (no Homebrew required):**
    - **Option 1: Go installer (recommended for Go users):**
      ```sh
      # Install the CLI
      go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
      # Verify installation
      migrate -version
      ```
    - **Option 2: Download prebuilt binary:**
      ```sh
      # Download the latest golang-migrate CLI for macOS (ARM64 example)
      curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.darwin-arm64.tar.gz | tar xvz
      # Or for Intel Macs, use migrate.darwin-amd64.tar.gz
      # For Linux, use migrate.linux-amd64.tar.gz, etc.
      mv migrate /usr/local/bin/
      chmod +x /usr/local/bin/migrate
      migrate -version
      ```
      > See [golang-migrate releases](https://github.com/golang-migrate/migrate/releases) for all OS/arch options.

2. **Create a new migration (replace `add_table` with your migration name):**
   ```sh
   migrate create -ext sql -dir migrations -seq add_table
   ```
   This will generate files like:
   ```
   migrations/000001_add_table.up.sql
   migrations/000001_add_table.down.sql
   ```
3. **Edit the generated `.up.sql` and `.down.sql` files with your migration SQL.**

> **Tip:**
> The `-seq` flag ensures sequential numbering. You can use any descriptive name after the number.
