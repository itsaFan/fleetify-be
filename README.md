# Fleetify Fullstack Challenge (Backend)

### Technology Used:

<p align="left">    
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/go/go-original-wordmark.svg" width="30"
                height="30" />
</p>

## Brief Description

#### Project Theme

Employee Attendances

#### Project Purpose:

Sebuah perusahaan Multinasional memiliki jumlah karyawan diatas 50 karyawan, dan memiliki berbagai macam Divisi atau departemen didalamnya. Karena banyaknya karyawan untuk dikelola, perusahaan membutuhkan Sistem untuk Absensi guna mencatat serta mengevaluasi kedisiplinan karyawan secara sistematis.

#### Guide to use this app on local

1. Git clone this repository.
2. Use `go mod tidy` on root folder to install all depedencies.
3. We have .env file, so you need to configure your own .env.
4. Initize your local MySQL database.
5. Refer to `.env.example` file on root directory.
6. You must run migration first after set the `GOOSE_DBSTRING` with `goose -env .env -dir db/migrations up`
7. Then to start the project on your local development `go run ./cmd/api`.

### Deployments

#### API-Documentation: [Postman](https://documenter.getpostman.com/view/43445325/2sB3HqGHzu)
