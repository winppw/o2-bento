# Frontend — Next.js 14

Stack: Next.js 14 App Router · TypeScript · Tailwind CSS

## Key files
- `src/lib/schedule.ts` — window open/close logic (single source of truth, mirrors backend)
- `src/components/OrderStatus.tsx` — live countdown, status badge, form button
- `src/app/page.tsx` — root page

## Rules
- All time logic goes through `getWindowStatus()` in `schedule.ts`
- Timezone is always Asia/Bangkok — never assume local or UTC
- The form button must be disabled and not link anywhere when `isOpen === false`
- Tailwind only for styling — no component libraries unless explicitly added
- `"use client"` only on components that use `useState`/`useEffect`
