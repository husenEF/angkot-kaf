import { Bot } from "grammy";
import { config } from "./config";
import { BotServiceImpl } from "./core/services/bot_service";
import { setupBot } from "./infrastructure/bot/telegram";
import { SQLiteDB } from "./infrastructure/database/sqlite";

async function main() {
    console.log("Starting bot...");

    try {
        const db = await SQLiteDB.initialize();
        const botService = new BotServiceImpl(db);
        const bot = new Bot(config.TELEGRAM_TOKEN);

        // Add middleware
        // bot.use(loggerMiddleware);

        // Setup bot commands
        setupBot(bot, botService);

        // Start bot
        await bot.api.setMyCommands([
            { command: "help", description: "Show available commands" },
            { command: "ping", description: "Check bot status" },
            { command: "passenger", description: "Add new passenger" },
            { command: "driver", description: "Add new driver" },
            { command: "passengers", description: "List all passengers" },
            { command: "drivers", description: "List all drivers" },
            { command: "report", description: "Get today's report" },
        ]);

        console.log("Bot is running...");
        await bot.start();
    } catch (error) {
        console.error("Failed to start bot:", error);
        process.exit(1);
    }
}

// Jalankan fungsi utama
main().catch((error) => {
    console.error("Error tidak tertangani:", error);
    process.exitCode = 1;
});