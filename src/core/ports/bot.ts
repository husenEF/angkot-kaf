export interface BotService {
    handlePing(): string;
    handlePassenger(chatId: number): string;
    addPassenger(name: string, chatId: number): Promise<void>;
    isWaitingForPassengerName(chatId: number): boolean;
    clearWaitingStatus(chatId: number): void;
    getPassengerList(chatId: number): Promise<string>;
    handleDriver(chatId: number): string;
    addDriver(name: string, chatId: number): Promise<void>;
    getDriverList(chatId: number): Promise<string>;
    isWaitingForDriverName(chatId: number): boolean;
    processDeparture(driverName: string, passengers: string[], chatId: number): Promise<string>;
    processReturn(driverName: string, passengers: string[], chatId: number): Promise<string>;
    getTodayReport(chatId: number): Promise<string>;
    getReportByDate(chatId: number, date: string): Promise<string>;
    parseAndProcessDeparture(text: string, chatId: number): Promise<string>;
    parseAndProcessReturn(text: string, chatId: number): Promise<string>;
}
