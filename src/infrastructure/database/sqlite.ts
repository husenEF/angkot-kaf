import { Database } from "sqlite3";
import type { Storage } from "../../core/ports/storage";

interface DriverRow {
    name: string;
}

interface PassengerRow {
    name: string;
}

interface TripRow {
    driver_name: string;
    passenger_name: string;
}

interface CountRow {
    count: number;
}

export class SQLiteDB implements Storage {
    private db: Database;

    private constructor(db: Database) {
        this.db = db;
    }

    static async initialize(): Promise<SQLiteDB> {
        const db = new Database("database/angkot.db");
        await this.initializeTables(db);
        return new SQLiteDB(db);
    }

    private static async initializeTables(db: Database): Promise<void> {
        const queries = [
            `CREATE TABLE IF NOT EXISTS drivers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        chat_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(name, chat_id)
      )`,
            `CREATE TABLE IF NOT EXISTS passengers (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        chat_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(name, chat_id)
      )`,
            `CREATE TABLE IF NOT EXISTS trips (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        driver_name TEXT NOT NULL,
        passenger_name TEXT NOT NULL,
        chat_id INTEGER NOT NULL,
        type TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
      )`
        ];

        for (const query of queries) {
            await new Promise<void>((resolve, reject) => {
                db.run(query, (err: Error | null) => {
                    if (err) reject(err);
                    else resolve();
                });
            });
        }
    }

    async saveDriver(name: string, chatId: number): Promise<void> {
        return new Promise<void>((resolve, reject) => {
            this.db.run(
                "INSERT INTO drivers (name, chat_id) VALUES (?, ?)",
                [name, chatId],
                (err: Error | null) => {
                    if (err) reject(err);
                    else resolve();
                }
            );
        });
    }

    async getDrivers(chatId: number): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<DriverRow>(
                "SELECT name FROM drivers WHERE chat_id = ? ORDER BY name",
                [chatId],
                (err: Error | null, rows: DriverRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.name));
                }
            );
        });
    }

    async isDriverExists(name: string, chatId: number): Promise<boolean> {
        return new Promise<boolean>((resolve, reject) => {
            this.db.get(
                "SELECT 1 FROM drivers WHERE name = ? AND chat_id = ?",
                [name, chatId],
                (err: Error | null, row: any) => {
                    if (err) reject(err);
                    else resolve(!!row);
                }
            );
        });
    }

    async savePassenger(name: string, chatId: number): Promise<void> {
        return new Promise<void>((resolve, reject) => {
            this.db.run(
                "INSERT INTO passengers (name, chat_id) VALUES (?, ?)",
                [name, chatId],
                (err: Error | null) => {
                    if (err) reject(err);
                    else resolve();
                }
            );
        });
    }

    async getPassengers(chatId: number): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<PassengerRow>(
                "SELECT name FROM passengers WHERE chat_id = ? ORDER BY name",
                [chatId],
                (err: Error | null, rows: PassengerRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.name));
                }
            );
        });
    }

    async saveDeparture(
        driverName: string,
        passengers: string[],
        chatId: number
    ): Promise<void> {
        const promises = passengers.map((passenger) => {
            return new Promise<void>((resolve, reject) => {
                this.db.run(
                    "INSERT INTO trips (driver_name, passenger_name, chat_id, type) VALUES (?, ?, ?, 'departure')",
                    [driverName, passenger, chatId],
                    (err: Error | null) => {
                        if (err) reject(err);
                        else resolve();
                    }
                );
            });
        });

        await Promise.all(promises);
    }

    async saveReturn(
        driverName: string,
        passengers: string[],
        chatId: number
    ): Promise<void> {
        const promises = passengers.map((passenger) => {
            return new Promise<void>((resolve, reject) => {
                this.db.run(
                    "INSERT INTO trips (driver_name, passenger_name, chat_id, type) VALUES (?, ?, ?, 'return')",
                    [driverName, passenger, chatId],
                    (err: Error | null) => {
                        if (err) reject(err);
                        else resolve();
                    }
                );
            });
        });

        await Promise.all(promises);
    }

    async getPassengerTripPrice(passengerName: string, chatId: number): Promise<number> {
        return new Promise<number>((resolve, reject) => {
            this.db.get<CountRow>(
                "SELECT COUNT(*) as count FROM trips WHERE passenger_name = ? AND chat_id = ? AND created_at >= date('now', 'start of month')",
                [passengerName, chatId],
                (err: Error | null, row: CountRow) => {
                    if (err) reject(err);
                    else resolve(row.count * 5000);
                }
            );
        });
    }

    async getDeparturePassengers(driverName: string, chatId: number): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<TripRow>(
                "SELECT DISTINCT passenger_name FROM trips WHERE driver_name = ? AND chat_id = ? AND type = 'departure' AND date(created_at) = date('now')",
                [driverName, chatId],
                (err: Error | null, rows: TripRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.passenger_name));
                }
            );
        });
    }

    async hasDepartureToday(passengerName: string, chatId: number): Promise<boolean> {
        return new Promise<boolean>((resolve, reject) => {
            this.db.get(
                "SELECT 1 FROM trips WHERE passenger_name = ? AND chat_id = ? AND type = 'departure' AND date(created_at) = date('now')",
                [passengerName, chatId],
                (err: Error | null, row: any) => {
                    if (err) reject(err);
                    else resolve(!!row);
                }
            );
        });
    }

    async getTodayDrivers(chatId: number): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<TripRow>(
                "SELECT DISTINCT driver_name FROM trips WHERE chat_id = ? AND date(created_at) = date('now')",
                [chatId],
                (err: Error | null, rows: TripRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.driver_name));
                }
            );
        });
    }

    async getReturnPassengers(driverName: string, chatId: number): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<TripRow>(
                "SELECT DISTINCT passenger_name FROM trips WHERE driver_name = ? AND chat_id = ? AND type = 'return' AND date(created_at) = date('now')",
                [driverName, chatId],
                (err: Error | null, rows: TripRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.passenger_name));
                }
            );
        });
    }

    async getDriversByDate(chatId: number, date: string): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<TripRow>(
                "SELECT DISTINCT driver_name FROM trips WHERE chat_id = ? AND date(created_at) = date(?)",
                [chatId, date],
                (err: Error | null, rows: TripRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.driver_name));
                }
            );
        });
    }

    async getDeparturePassengersByDate(driverName: string, chatId: number, date: string): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<TripRow>(
                "SELECT DISTINCT passenger_name FROM trips WHERE driver_name = ? AND chat_id = ? AND type = 'departure' AND date(created_at) = date(?)",
                [driverName, chatId, date],
                (err: Error | null, rows: TripRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.passenger_name));
                }
            );
        });
    }

    async getReturnPassengersByDate(driverName: string, chatId: number, date: string): Promise<string[]> {
        return new Promise<string[]>((resolve, reject) => {
            this.db.all<TripRow>(
                "SELECT DISTINCT passenger_name FROM trips WHERE driver_name = ? AND chat_id = ? AND type = 'return' AND date(created_at) = date(?)",
                [driverName, chatId, date],
                (err: Error | null, rows: TripRow[]) => {
                    if (err) reject(err);
                    else resolve(rows.map(row => row.passenger_name));
                }
            );
        });
    }
}