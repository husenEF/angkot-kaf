import { Bot, Context } from "grammy";
import type { BotService } from "../../core/ports/bot";

const HELP_MESSAGE = `
Available commands:
/ping - Check bot status
/passenger - Add new passenger
/driver - Add new driver
/passengers - List all passengers
/drivers - List all drivers
/departure <driver> - <passenger1>, <passenger2> - Record departure
/return <driver> - <passenger1>, <passenger2> - Record return
/report - Get today's report
/report_date YYYY-MM-DD - Get report by date
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

        if (!message) return;

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
    });
}