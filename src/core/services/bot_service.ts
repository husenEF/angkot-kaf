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

        await this.storage.saveDeparture(driverName, passengers, chatId);
        return "Keberangkatan berhasil dicatat";
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

        await this.storage.saveReturn(driverName, passengers, chatId);
        return "Kepulangan berhasil dicatat";
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
}