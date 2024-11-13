import { ROUND_TRIP_PRICE, SINGLE_TRIP_PRICE } from "../constants/price";
import type { BotService } from "../ports/bot";
import type { Storage } from "../ports/storage";

export class BotServiceImpl implements BotService {
    private storage: Storage;
    private waitingForPassengerName: Map<number, boolean>;
    private waitingForDriverName: Map<number, boolean>;

    constructor(storage: Storage) {
        this.storage = storage;
        this.waitingForPassengerName = new Map();
        this.waitingForDriverName = new Map();
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

    async processDeparture(
        driverName: string,
        passengers: string[],
        chatId: number
    ): Promise<string> {
        const driverExists = await this.storage.isDriverExists(driverName, chatId);
        if (!driverExists) {
            return `Supir ${driverName} tidak terdaftar`;
        }

        const totalAmount = SINGLE_TRIP_PRICE * passengers.length;
        let priceDetails = "Detail Pembayaran:\n";
        passengers.forEach(passenger => {
            priceDetails += `${passenger}: Rp ${SINGLE_TRIP_PRICE.toLocaleString('id-ID')}\n`;
        });

        await this.storage.saveDeparture(driverName, passengers, chatId);

        return `Keberangkatan berhasil dicatat\n\n${priceDetails}\nTotal untuk driver: Rp ${totalAmount.toLocaleString('id-ID')}`;
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

        // Calculate prices for each passenger
        let totalAmount = 0;
        let priceDetails = "Detail Pembayaran:\n";

        for (const passenger of passengers) {
            const hasDepartureToday = await this.storage.hasDepartureToday(passenger, chatId);
            const price = hasDepartureToday ? ROUND_TRIP_PRICE - SINGLE_TRIP_PRICE : SINGLE_TRIP_PRICE;
            totalAmount += price;
            priceDetails += `${passenger}: Rp ${price.toLocaleString('id-ID')}\n`;
        }

        await this.storage.saveReturn(driverName, passengers, chatId);

        return `Kepulangan berhasil dicatat\n\n${priceDetails}\nTotal untuk driver: Rp ${totalAmount.toLocaleString('id-ID')}`;
    }

    async getTodayReport(chatId: number): Promise<string> {
        const drivers = await this.storage.getTodayDrivers(chatId);
        if (drivers.length === 0) {
            return "Belum ada perjalanan hari ini";
        }

        let report = "Laporan Hari Ini:\n\n";
        for (const driver of drivers) {
            const departurePassengers = await this.storage.getDeparturePassengers(
                driver,
                chatId
            );
            const returnPassengers = await this.storage.getReturnPassengers(
                driver,
                chatId
            );

            report += `Supir: ${driver}\n`;
            report += "Berangkat: " + departurePassengers.join(", ") + "\n";
            report += "Pulang: " + returnPassengers.join(", ") + "\n\n";
        }

        return report;
    }

    async getReportByDate(chatId: number, date: string): Promise<string> {
        const drivers = await this.storage.getDriversByDate(chatId, date);
        if (drivers.length === 0) {
            return `Belum ada perjalanan pada tanggal ${date}`;
        }

        let report = `Laporan Tanggal ${date}:\n\n`;
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

            report += `Supir: ${driver}\n`;
            report += "Berangkat: " + departurePassengers.join(", ") + "\n";
            report += "Pulang: " + returnPassengers.join(", ") + "\n\n";
        }

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
}