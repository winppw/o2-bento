export interface WindowStatus {
  isOpen: boolean;
  openAt: Date;
  closeAt: Date;
  isFridayExtra: boolean;
  message: string;
}

const BKK_OFFSET_MS = 7 * 60 * 60 * 1000; // UTC+7

export function getWindowStatus(now: Date = new Date()): WindowStatus {
  // Shift now by +7 h so that getUTC* methods read Bangkok wall-clock values.
  const proxy = new Date(now.getTime() + BKK_OFFSET_MS);
  const bkkHour = proxy.getUTCHours();
  const bkkMinute = proxy.getUTCMinutes();
  const bkkDay = proxy.getUTCDay();
  const y = proxy.getUTCFullYear();
  const mo = proxy.getUTCMonth();
  const d = proxy.getUTCDate();

  // Real UTC timestamps for 11:00 and 16:00 Bangkok time on the same calendar day.
  const openAt = new Date(Date.UTC(y, mo, d, 11 - 7, 0, 0, 0));
  const closeAt = new Date(Date.UTC(y, mo, d, 16 - 7, 0, 0, 0));

  const isWeekend = bkkDay === 0 || bkkDay === 6;
  const bkkMins = bkkHour * 60 + bkkMinute;
  const isOpen = !isWeekend && bkkMins >= 11 * 60 && bkkMins < 16 * 60;
  const isFriday = bkkDay === 5;

  let message = "";
  if (!isWeekend && bkkMins < 11 * 60) {
    message = "Order window opens at 11:00 AM.";
  } else if (isOpen && isFriday) {
    message = "Order window is open. You can also submit Monday's order today.";
  } else if (isOpen) {
    message = "Order window is open. Deadline is 4:00 PM.";
  } else {
    message = "Order window is closed for today.";
  }

  return { isOpen, openAt, closeAt, isFridayExtra: isFriday, message };
}

export function formatTimeRemaining(target: Date, now: Date = new Date()): string {
  const diff = target.getTime() - now.getTime();
  if (diff <= 0) return "0:00";
  const h = Math.floor(diff / 3_600_000);
  const m = Math.floor((diff % 3_600_000) / 60_000);
  const s = Math.floor((diff % 60_000) / 1_000);
  if (h > 0) return `${h}h ${m}m`;
  return `${m}:${s.toString().padStart(2, "0")}`;
}
