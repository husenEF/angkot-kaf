import { Bot, Context } from "grammy";
import type { BotService } from "../../core/ports/bot";

const HELP_MESSAGE = `
Perintah yang tersedia:
/ping - Cek status bot
/passenger - Tambah penumpang baru
/driver - Tambah supir baru
/passengers - Daftar semua penumpang
/drivers - Daftar semua supir
/report - Laporan hari ini
/report_date YYYY-MM-DD - Laporan per tanggal

Format input perjalanan:
antar
Driver: [nama_supir]
- [penumpang1]
- [penumpang2]

jemput
Driver: [nama_supir]
- [penumpang1]
- [penumpang2]
`;

export function setupBot(bot: Bot, service: BotService): void {
    bot.command("ping", (ctx: Context) => ctx.reply(service.handlePing()));

    bot.command("passenger", async (ctx: Context) => {
        const chatId: number = ctx.chat?.id ?? 0;
        const response: string = service.handlePassenger(chatId);
        await ctx.reply(response);
    });

    bot.command("driver", async (ctx: Context) => {
        const chatId: number = ctx.chat?.id ?? 0;
        const response: string = service.handleDriver(chatId);
        await ctx.reply(response);
    });

    bot.command("passengers", async (ctx: Context) => {
        const chatId: number = ctx.chat?.id ?? 0;
        const response: string = await service.getPassengerList(chatId);
        await ctx.reply(response);
    });

    bot.command("drivers", async (ctx: Context) => {
        const chatId: number = ctx.chat?.id ?? 0;
        const response: string = await service.getDriverList(chatId);
        await ctx.reply(response);
    });

    bot.command("report", async (ctx: Context) => {
        const chatId: number = ctx.chat?.id ?? 0;
        const response: string = await service.getTodayReport(chatId);
        await ctx.reply(response);
    });

    bot.on("message", async (ctx: Context) => {
        const chatId: number = ctx.chat?.id ?? 0;
        const message: string | undefined = ctx.message?.text;

        console.log({
            timestamp: new Date().toISOString(),
            chatId: chatId,
            message: message,
            username: ctx.from?.username,
            firstName: ctx.from?.first_name
        });

        if (!message) return;

        if (message.toLowerCase().startsWith('antar')) {
            const inputText = message.substring(5).trim();
            const result = await service.parseAndProcessTrip(inputText, chatId, 'antar');
            await ctx.reply(result);
            return;
        }

        if (message.toLowerCase().startsWith('jemput')) {
            const inputText = message.substring(6).trim();
            const result = await service.parseAndProcessTrip(inputText, chatId, 'jemput');
            await ctx.reply(result);
            return;
        }

        if (service.isWaitingForPassengerName(chatId)) {
            await service.addPassenger(message, chatId);
            service.clearWaitingStatus(chatId);
            await ctx.reply(`Penumpang ${message} berhasil ditambahkan`);
            return;
        }

        if (service.isWaitingForDriverName(chatId)) {
            await service.addDriver(message, chatId);
            service.clearWaitingStatus(chatId);
            await ctx.reply(`Supir ${message} berhasil ditambahkan`);
            return;
        }

        if (message.toLowerCase() === '/help' || message.toLowerCase() === 'help') {
            await ctx.reply(HELP_MESSAGE);
            return;
        }
    });
}