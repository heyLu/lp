import { LiveReload, Outlet } from "@remix-run/react";

export default function App() {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <title>quoth me, she said</title>
      </head>

      <body>
        <Outlet />
        {/*<Scripts />*/}
        <LiveReload />
      </body>
    </html>
  );
}
