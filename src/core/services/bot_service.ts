import { ROUND_TRIP_PRICE, SINGLE_TRIP_PRICE } from "../constants/price";
import type { BotService } from "../ports/bot";
import type { Storage } from "../ports/storage";

export class BotServiceImpl implements BotService {
    private storage: Storage;
    private waitingForPassengerName: Map<number, boolean>;
    private waitingForDriverName: Map<number, boolean>;
    private waitingForCatatan: Map<number, boolean>;

    constructor(storage: Storage) {
        this.storage = storage;
        this.waitingForPassengerName = new Map();
        this.waitingForDriverName = new Map();
        this.waitingForCatatan = new Map();
    }

    handlePing(): string {
        return "pong";
    }

    handlePassenger(chatId: number): string {
        this.waitingForPassengerName.set(chatId, true);
        return "Masukkan nama penumpang:";
    }

    async addPassenger(name: string, chatId: number): Promise<void> {
        await this.storage.savePassenger(name, chatId);
    }

    isWaitingForPassengerName(chatId: number): boolean {
        return this.waitingForPassengerName.get(chatId) || false;
    }

    clearWaitingStatus(chatId: number): void {
        this.waitingForPassengerName.delete(chatId);
        this.waitingForDriverName.delete(chatId);
        this.waitingForCatatan.delete(chatId);
    }

    async getPassengerList(chatId: number): Promise<string> {
        const passengers = await this.storage.getPassengers(chatId);
        if (passengers.length === 0) {
            return "Belum ada penumpang yang terdaftar";
        }
        return "Daftar Penumpang:\n" + passengers.join("\n");
    }

    handleDriver(chatId: number): string {
        this.waitingForDriverName.set(chatId, true);
        return "Masukkan nama supir:";
    }

    async addDriver(name: string, chatId: number): Promise<void> {
        await this.storage.saveDriver(name, chatId);
    }

    async getDriverList(chatId: number): Promise<string> {
        const drivers = await this.storage.getDrivers(chatId);
        if (drivers.length === 0) {
            return "Belum ada supir yang terdaftar";
        }
        return "Daftar Supir:\n" + drivers.join("\n");
    }

    isWaitingForDriverName(chatId: number): boolean {
        return this.waitingForDriverName.get(chatId) || false;
    }

    isWaitingForCatatan(chatId: number): boolean {
        return this.waitingForCatatan.get(chatId) || false;
    }

    async processDeparture(
        driverName: string,
        passengers: string[],
        chatId: number
    ): Promise<string> {
        const driverExists = await this.storage.isDriverExists(driverName, chatId);
        if (!driverExists) {
            return `Supir ${driverName} tidak terdaftar`;
        }

        // Hapus data departure sebelumnya untuk driver ini pada hari yang sama
        await this.storage.deleteDepartureToday(driverName, chatId);

        // Simpan data departure yang baru
        await this.storage.saveDeparture(driverName, passengers, chatId);

        return `Keberangkatan berhasil dicatat\nDriver: ${driverName}\nPenumpang:\n${passengers.map(p => `- ${p}`).join('\n')}`;
    }

    async processReturn(
        driverName: string,
        passengers: string[],
        chatId: number
    ): Promise<string> {
        const driverExists = await this.storage.isDriverExists(driverName, chatId);
        if (!driverExists) {
            return `Supir ${driverName} tidak terdaftar`;
        }

        // Hapus data return sebelumnya untuk driver ini pada hari yang sama
        await this.storage.deleteReturnToday(driverName, chatId);

        // Dapatkan daftar penumpang yang berangkat hari ini
        const departurePassengers = await this.storage.getDeparturePassengers(driverName, chatId);

        // Identifikasi penumpang yang berangkat tapi tidak pulang
        const notReturningPassengers = departurePassengers.filter(p => !passengers.includes(p));

        let response = `Kepulangan berhasil dicatat\nDriver: ${driverName}\nPenumpang:\n${passengers.map(p => `- ${p}`).join('\n')}`;

        // Tambahkan informasi penumpang yang tidak ikut pulang
        if (notReturningPassengers.length > 0) {
            response += `\n\nTidak Ikut Pulang:\n${notReturningPassengers.map(p => `- ${p}`).join('\n')}`;
        }

        // Simpan data return yang baru
        await this.storage.saveReturn(driverName, passengers, chatId);

        return response;
    }

    async getTodayReport(chatId: number): Promise<string> {
        const drivers = await this.storage.getTodayDrivers(chatId);
        if (drivers.length === 0) {
            return "Belum ada perjalanan hari ini";
        }

        let report = "Laporan\n\n";
        let totalAllDrivers = 0;

        for (const driver of drivers) {
            const departurePassengers = await this.storage.getDeparturePassengers(driver, chatId);
            const returnPassengers = await this.storage.getReturnPassengers(driver, chatId);

            let driverTotal = 0;
            let roundTripPassengers = [];
            let departureOnlyPassengers = [];
            let returnOnlyPassengers = [];

            // Kategorikan penumpang
            for (const passenger of departurePassengers) {
                if (returnPassengers.includes(passenger)) {
                    roundTripPassengers.push(passenger);
                    driverTotal += ROUND_TRIP_PRICE;
                } else {
                    departureOnlyPassengers.push(passenger);
                    driverTotal += SINGLE_TRIP_PRICE;
                }
            }

            // Cari penumpang yang hanya pulang
            returnOnlyPassengers = returnPassengers.filter(p => !departurePassengers.includes(p));
            driverTotal += returnOnlyPassengers.length * SINGLE_TRIP_PRICE;

            report += `Driver: ${driver}\n\n`;

            // Pulang Pergi
            if (roundTripPassengers.length > 0) {
                report += "Angkutan Pulang Pergi\n";
                roundTripPassengers.forEach(passenger => {
                    report += `- ${passenger} (Rp ${ROUND_TRIP_PRICE.toLocaleString('id-ID')})\n`;
                });
                report += "\n";
            }

            // Berangkat Saja
            if (departureOnlyPassengers.length > 0) {
                report += "Berangkat Saja\n";
                departureOnlyPassengers.forEach(passenger => {
                    report += `- ${passenger} (Rp ${SINGLE_TRIP_PRICE.toLocaleString('id-ID')})\n`;
                });
                report += "\n";
            }

            // Pulang Saja
            if (returnOnlyPassengers.length > 0) {
                report += "Pulang Saja\n";
                returnOnlyPassengers.forEach(passenger => {
                    report += `- ${passenger} (Rp ${SINGLE_TRIP_PRICE.toLocaleString('id-ID')})\n`;
                });
                report += "\n";
            }

            report += `Total Untuk ${driver}: Rp ${driverTotal.toLocaleString('id-ID')}\n\n`;
            totalAllDrivers += driverTotal;
        }

        report += `Total Keseluruhan: Rp ${totalAllDrivers.toLocaleString('id-ID')}`;
        return report;
    }

    async getReportByDate(chatId: number, date: string): Promise<string> {
        const drivers = await this.storage.getDriversByDate(chatId, date);
        if (drivers.length === 0) {
            return `Belum ada perjalanan pada tanggal ${date}`;
        }

        let report = `Laporan Tanggal ${date}:\n\n`;
        let totalAllDrivers = 0;

        for (const driver of drivers) {
            const departurePassengers = await this.storage.getDeparturePassengersByDate(
                driver,
                chatId,
                date
            );
            const returnPassengers = await this.storage.getReturnPassengersByDate(
                driver,
                chatId,
                date
            );

            let driverTotal = 0;
            let tripDetails = "";

            // Hitung biaya keberangkatan
            if (departurePassengers.length > 0) {
                const departureTotal = SINGLE_TRIP_PRICE * departurePassengers.length;
                driverTotal += departureTotal;
                tripDetails += "Berangkat:\n";
                departurePassengers.forEach(passenger => {
                    tripDetails += `- ${passenger}: Rp ${SINGLE_TRIP_PRICE.toLocaleString('id-ID')}\n`;
                });
            }

            // Hitung biaya kepulangan
            if (returnPassengers.length > 0) {
                tripDetails += "\nPulang:\n";
                for (const passenger of returnPassengers) {
                    const hasDepartureToday = departurePassengers.includes(passenger);
                    const returnPrice = hasDepartureToday ?
                        ROUND_TRIP_PRICE - SINGLE_TRIP_PRICE :
                        SINGLE_TRIP_PRICE;
                    driverTotal += returnPrice;
                    tripDetails += `- ${passenger}: Rp ${returnPrice.toLocaleString('id-ID')}${hasDepartureToday ? ' (Pulang PP)' : ''}\n`;
                }
            }

            // Tambahkan ke total keseluruhan
            totalAllDrivers += driverTotal;

            report += `Supir: ${driver}\n`;
            report += tripDetails;
            report += `Total untuk ${driver}: Rp ${driverTotal.toLocaleString('id-ID')}\n\n`;
        }

        report += `Total keseluruhan: Rp ${totalAllDrivers.toLocaleString('id-ID')}`;
        return report;
    }

    async parseAndProcessTrip(text: string, chatId: number, type: 'antar' | 'jemput'): Promise<string> {
        try {
            const lines = text.trim().split('\n');

            const driverLine = lines[0];
            const driverMatch = driverLine.match(/Driver:\s*(.+)/i);
            if (!driverMatch) {
                return "Format salah. Gunakan format:\nDriver: [nama]\n- [penumpang1]\n- [penumpang2]";
            }
            const driverName = driverMatch[1].trim();

            const passengers: string[] = [];
            for (let i = 1; i < lines.length; i++) {
                const line = lines[i].trim();
                if (line.startsWith('-')) {
                    const passengerName = line.substring(1).trim();
                    if (passengerName) {
                        passengers.push(passengerName);
                    }
                }
            }

            if (passengers.length === 0) {
                return "Tidak ada penumpang yang tercantum";
            }

            // Memproses sesuai tipe (antar atau jemput)
            if (type === 'antar') {
                return await this.processDeparture(driverName, passengers, chatId);
            } else {
                return await this.processReturn(driverName, passengers, chatId);
            }
        } catch (error) {
            console.error(`Error parsing ${type} text:`, error);
            return "Terjadi kesalahan saat memproses input. Pastikan format sudah benar.";
        }
    }

    async backupDatabase(chatId: number): Promise<{ path: string; filename: string }> {
        const dbPath = "database/angkot.db";
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const backupFilename = `angkot-backup-${timestamp}.db`;

        return {
            path: dbPath,
            filename: backupFilename
        };
    }

    handleCatat(chatId: number): string {
        this.waitingForCatatan.set(chatId, true);
        return "Silakan masukkan catatan perjalanan dengan format:\n\n" +
            "Hari, Tanggal\n" +
            "Driver: [nama_driver]\n\n" +
            "Antar & Jemput:\n" +
            "1. [nama_penumpang]\n" +
            "2. [nama_penumpang]\n\n" +
            "Antar aja:\n" +
            "1. [nama_penumpang]\n\n" +
            "Jemput aja:\n" +
            "1. [nama_penumpang]";
    }

    async processCatatanPerjalanan(text: string, chatId: number): Promise<string> {
        try {
            const lines = text.split('\n').map(line => line.trim());
            const driverMatch = lines.find(line => line.toLowerCase().startsWith('driver'))?.match(/driver\s*:\s*(.+)/i);

            if (!driverMatch) {
                return "Format salah. Mohon sertakan nama driver dengan format 'Driver: [nama]'";
            }

            const driverName = driverMatch[1].trim();
            const roundTripPassengers: string[] = [];
            const departureOnlyPassengers: string[] = [];
            const returnOnlyPassengers: string[] = [];

            let currentSection = '';

            for (const line of lines) {
                if (line.toLowerCase().includes('antar & jemput')) {
                    currentSection = 'roundtrip';
                    continue;
                } else if (line.toLowerCase().includes('antar aja')) {
                    currentSection = 'departure';
                    continue;
                } else if (line.toLowerCase().includes('jemput aja')) {
                    currentSection = 'return';
                    continue;
                }

                const passengerMatch = line.match(/^\d+\.\s*(.+)$/);
                if (passengerMatch) {
                    const passengerName = passengerMatch[1].trim();
                    switch (currentSection) {
                        case 'roundtrip':
                            roundTripPassengers.push(passengerName);
                            break;
                        case 'departure':
                            departureOnlyPassengers.push(passengerName);
                            break;
                        case 'return':
                            returnOnlyPassengers.push(passengerName);
                            break;
                    }
                }
            }

            // Simpan data perjalanan
            if (roundTripPassengers.length > 0) {
                await this.storage.saveDeparture(driverName, roundTripPassengers, chatId);
                await this.storage.saveReturn(driverName, roundTripPassengers, chatId);
            }

            if (departureOnlyPassengers.length > 0) {
                await this.storage.saveDeparture(driverName, departureOnlyPassengers, chatId);
            }

            if (returnOnlyPassengers.length > 0) {
                await this.storage.saveReturn(driverName, returnOnlyPassengers, chatId);
            }

            // Buat laporan
            let report = `✅ Catatan berhasil disimpan\n\n`;
            report += `Driver: ${driverName}\n\n`;

            if (roundTripPassengers.length > 0) {
                report += `Antar & Jemput:\n${roundTripPassengers.map((p, i) => `${i + 1}. ${p}`).join('\n')}\n\n`;
            }

            if (departureOnlyPassengers.length > 0) {
                report += `Antar saja:\n${departureOnlyPassengers.map((p, i) => `${i + 1}. ${p}`).join('\n')}\n\n`;
            }

            if (returnOnlyPassengers.length > 0) {
                report += `Jemput saja:\n${returnOnlyPassengers.map((p, i) => `${i + 1}. ${p}`).join('\n')}`;
            }

            return report;
        } catch (error) {
            console.error('Error processing catatan:', error);
            return "Terjadi kesalahan saat memproses catatan. Pastikan format sudah benar.";
        }
    }
}