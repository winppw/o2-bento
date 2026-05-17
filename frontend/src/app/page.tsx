import OrderStatus from "@/components/OrderStatus";

export default function Home() {
  return (
    <main className="w-full max-w-lg p-4">
      <h1 className="text-2xl font-bold text-center text-gray-800 mb-2">🍱 Lunch Order</h1>
      <p className="text-center text-gray-500 text-sm mb-8">
        Submit daily via Google Form · Deadline 4:00 PM
      </p>
      <OrderStatus />
    </main>
  );
}
