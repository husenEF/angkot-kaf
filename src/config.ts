import { z } from "zod";

// Definisikan schema untuk validasi environment variables
const envSchema = z.object({
    TELEGRAM_TOKEN: z.string({
        required_error: "TELEGRAM_TOKEN harus diisi",
    }),
    ADMIN_ID: z.string({
        required_error: "ADMIN_ID harus diisi",
    }).transform(Number),
});

// Parse dan validasi environment variables langsung dari process.env
const parsed = envSchema.safeParse({
    TELEGRAM_TOKEN: process.env.TELEGRAM_TOKEN,
    ADMIN_ID: process.env.ADMIN_ID
});

if (!parsed.success) {
    console.error("❌ Validasi environment variables gagal:", parsed.error.toString());
    process.exit(1);
}

export const config = parsed.data;