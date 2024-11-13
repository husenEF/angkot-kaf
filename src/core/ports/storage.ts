export interface Storage {
    saveDriver(name: string, chatId: number): Promise<void>;
    getDrivers(chatId: number): Promise<string[]>;
    isDriverExists(name: string, chatId: number): Promise<boolean>;
    savePassenger(name: string, chatId: number): Promise<void>;
    getPassengers(chatId: number): Promise<string[]>;
    saveDeparture(driverName: string, passengers: string[], chatId: number): Promise<void>;
    saveReturn(driverName: string, passengers: string[], chatId: number): Promise<void>;
    getPassengerTripPrice(passengerName: string, chatId: number): Promise<number>;
    getDeparturePassengers(driverName: string, chatId: number): Promise<string[]>;
    hasDepartureToday(passengerName: string, chatId: number): Promise<boolean>;
    getTodayDrivers(chatId: number): Promise<string[]>;
    getReturnPassengers(driverName: string, chatId: number): Promise<string[]>;
    getDriversByDate(chatId: number, date: string): Promise<string[]>;
    getDeparturePassengersByDate(driverName: string, chatId: number, date: string): Promise<string[]>;
    getReturnPassengersByDate(driverName: string, chatId: number, date: string): Promise<string[]>;
}