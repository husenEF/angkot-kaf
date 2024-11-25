# Angkot

Bot Telegram yang membantu untuk mengelola data angkot KAF

## Cara penggunaan
1. /start

# Fitur
1. /add_driver - Melakukan Pendaftaran Driver dengan push tombol add driver dari menu utama
2. /gas - Melakukan Pencatatan Antar Jemput dari menu utama 
3. /report - Melihar catatan record hari ini
4. /backupdb - Melakukan Backup Database

## Biaya Perjalanan
Biaya Perjalanan untuk santri yang ikut berangkat saja adalah 10.000 rupiah
Biaya perjalanan untuk santri yang ikut antar dan jemput adalah 15.000 rupiah
Driver akan mendapatkan total biaya perjalanan sesuai dengan jumlah santri yang ikut

## Format 
Format ini adalah contoh pencatatan dengan perjalanan antar jemput

### Format All

Rabu, 13 November 2024
Driver : Wisnu

1. Alice
2. Syafiq  
3. ⁠jibs
4. ⁠Mumtaza
5. Omar
6. asya
7. rasqa

Antar saja
1. Hamza
2. Mikayla
3. Ibrahim

### Format Antar 

antar
driver: Wisnu
1. Syafiq
2. ⁠Mumtaza
4. Omar

### format jemput 

jemput
driver: Wisnu
1. Syafiq
2. ⁠Mumtaza
4. Omar

## Rules 

- Driver harus terdaftar terlebih dahulu
- Penumpang tidak harus terdaftar 
- Ketika melakukan replace, misal saya sudah mengirimkan format antar, lalu saya bisa me-replace dengan format antar yang baru, maka format antar yang lama akan di hapus dan replace 
- Ketika melakukan replace, misal saya sudah mengirimkan format jemput, lalu saya bisa me-replace dengan format jemput yang baru, maka format jemput yang lama akan di hapus dan replace
- Format laporan buat dengan ada emoji supaya lebih menarik

## Stack
- Golang
- GORM 
- SQLite di folder (database/angkot.db)