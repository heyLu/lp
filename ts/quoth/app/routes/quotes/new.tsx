import { useState } from "react";

function toEmbedUrl(url: string): string {
  try {
    const u = new URL(url);
    if (u.hostname.includes("youtu") && u.searchParams.has("v")) {
      return `https://www.youtube.com/embed/${u.searchParams.get("v")}`;
    }
  } catch {
    // fall through
  }
  return url;
}

export default function NewQuoteRoute() {
  const [videoUrl, setVideoUrl] = useState("");

  return (
    <div>
      <form method="post">
        <div>
          <label>
            text: <textarea name="content" placeholder="what was said" cols={100} rows={5} required></textarea>
          </label>
        </div>
        <div>
          <label>
            by: <input name="author" type="text" placeholder="author" required />
          </label>
        </div>
        <div>
          <label>
            quote source: <input name="source" type="url" placeholder="link to source" size={50} required onInput={(ev) => ev.target.value.length > 10 && setVideoUrl(new URL(ev.target.value).toString())} />
          </label>
        </div>
        {videoUrl.length > 10 &&
          <div>
            <iframe src={toEmbedUrl(videoUrl)} width={400} height={200} />
          </div>
        }
        <div>
          <label>
            start: <input name="start" type="number" placeholder="in seconds" required />
          </label>
        </div>
        <div>
          <label>
            end: <input name="end" type="number" placeholder="in seconds" required />
          </label>
        </div>
        <div>
          <button type="submit" className="button">
            quote it!
          </button>
        </div>
      </form>
    </div>
  );
}
