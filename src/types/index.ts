export interface Trip {
    driverName: string;
    passengerName: string;
    type: 'departure' | 'return';
    createdAt: Date;
}

export interface Report {
    date: string;
    trips: Trip[];
}