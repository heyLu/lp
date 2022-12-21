import { Outlet } from "@remix-run/react";

export default function QuotesRoute() {
  return (
    <div>
      <h1>quotes! 💬</h1>
      <main>
        <Outlet />
      </main>
    </div>
  );
}
