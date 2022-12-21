function toEmbedUrl(quote): string | undefined {
  try {
    const u = new URL(quote.url);
    if (u.hostname.includes("youtu") && u.searchParams.has("v")) {
      return `https://www.youtube.com/embed/${u.searchParams.get("v")}?start=${quote.start}&end=${quote.end}&loop=1&controls=0`;
    }
  } catch {
    // fall through
  }
  return undefined;
}

export default function QuoteRoute() {
  const quote = {
    id: "mocked",
    text: "and they are mocked, teased, ribbed, mercilessly by their colleagues for being a little bit weird about coffee as if that's a bad thing.",
    author: "James Hoffmann",
    url: "https://www.youtube.com/watch?v=4xOEIpbxM4w",
    start: 125,
    end: 133,
    quoter: "ljunes",
  }

  const embedUrl = toEmbedUrl(quote);

  return (
    <article>
      <h1>mocked</h1>

      {(embedUrl &&
        <iframe src={embedUrl} width={400} height={200} />)
        || <a href={quote.url}>{quote.url}</a>}

      <p><q cite={quote.url}>{quote.text}</q> — {quote.author}</p>

      <footer>
        <p>found by {quote.quoter}</p>
      </footer>
    </article>
  );
}
