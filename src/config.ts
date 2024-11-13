import * as dotenv from "dotenv";
import { z } from "zod";

// Load environment variables dari file .env
dotenv.config();

// Definisikan schema untuk validasi environment variables
const envSchema = z.object({
    TELEGRAM_TOKEN: z.string({
        required_error: "TELEGRAM_TOKEN harus diisi",
    }),
    ADMIN_ID: z.string({
        required_error: "ADMIN_ID harus diisi",
    }).transform(Number),
});

// Parse dan validasi environment variables
const parsed = envSchema.safeParse(process.env);

if (!parsed.success) {
    console.error("❌ Validasi environment variables gagal:", parsed.error.toString());
    process.exit(1);
}

// Export config yang sudah divalidasi
export const config = parsed.data;