# Angkot KAF 🚐

Bot Telegram untuk mengelola data transportasi dan catatan pengemudi KAF. Bot ini membantu melacak layanan antar-jemput santri secara efisien.

## Fitur 🌟

- **Manajemen Pengemudi** (`/add_driver`)
  - Mendaftarkan pengemudi baru
  - Mengelola data pengemudi

- **Pencatatan Perjalanan** (`/gas`)
  - Mencatat perjalanan antar-jemput santri
  - Melacak perjalanan satu arah dan pulang-pergi
  - Perhitungan biaya otomatis

- **Pelaporan** (`/report`)
  - Melihat catatan perjalanan hari ini
  - Melacak riwayat transportasi santri

- **Backup Database** (`/backupdb`)
  - Backup database sistem
  - Menjaga keamanan data

## Struktur Biaya 💰

- Perjalanan satu arah (antar ATAU jemput): Rp. 10.000 per santri
- Perjalanan pulang-pergi (antar DAN jemput): Rp. 15.000 per santri
- Pengemudi menerima total biaya berdasarkan jumlah santri yang dilayani

## Instalasi 🔧

1. Clone repository:
   ```bash
   git clone https://github.com/robzlabz/angkot-kaf.git
   cd angkot-kaf
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Pengaturan environment variables:
   - Salin `.env.example` ke `.env`
   - Konfigurasi Token Bot Telegram Anda
   ```bash
   cp .env.example .env
   ```

4. Build dan jalankan:
   ```bash
   make build
   ./angkot
   ```

## Panduan Penggunaan 📝

### Perintah Dasar
1. `/start` - Memulai bot
2. `/add_driver` - Mendaftarkan pengemudi baru
3. `/gas` - Mencatat perjalanan baru
4. `/report` - Melihat catatan hari ini
5. `/backupdb` - Backup database

### Contoh Format Pencatatan

#### Format Antar-Jemput
```
Rabu, 13 November 2024
Driver: Wisnu

1. Alice
2. Syafiq  
3. Jibs
4. Mumtaza
5. Omar
6. Asya
7. Rasqa

Antar saja:
1. Hamza
2. Mikayla
3. Ibrahim
```

#### Format Antar
```
antar
driver: Wisnu
1. Syafiq
2. Mumtaza
3. Omar
```

#### Format Jemput
```
jemput
driver: Wisnu
1. Syafiq
2. Mumtaza
3. Omar
```

## Stack Teknis 🛠

- **Bahasa**: Go (1.22.3)
- **Database**: SQLite (terletak di `database/angkot.db`)
- **Dependencies**:
  - `github.com/go-telegram-bot-api/telegram-bot-api/v5`: Telegram Bot API
  - `gorm.io/gorm`: ORM untuk operasi database
  - `gorm.io/driver/sqlite`: Driver SQLite
  - `github.com/joho/godotenv`: Konfigurasi environment

## Aturan dan Panduan 📋

1. Pengemudi harus terdaftar sebelum mencatat perjalanan
2. Penumpang tidak perlu terdaftar terlebih dahulu
3. Format perjalanan dapat diganti dengan yang baru (catatan sebelumnya akan diperbarui)
4. Laporan mencakup emoji untuk keterbacaan yang lebih baik
5. Format laporan dibuat dengan emoji agar lebih menarik

## Cara Berkontribusi 🤝

1. Fork repository ini
2. Buat branch fitur Anda (`git checkout -b fitur/fitur-keren`)
3. Commit perubahan Anda (`git commit -m 'Menambahkan fitur keren'`)
4. Push ke branch (`git push origin fitur/fitur-keren`)
5. Buat Pull Request

## Lisensi 📄

Proyek ini dilisensikan di bawah Lisensi MIT - lihat file LICENSE untuk detailnya.
