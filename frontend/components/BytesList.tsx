import { DateTime } from "luxon"
import { FC } from "react"
import { bytes } from "../client/client"

const BytesList: FC<{ bytes: bytes.Byte[] }> = ({ bytes }) => (
  <>
    {bytes.map((byte) => {
      const created = DateTime.fromISO(byte.created)
      return (
        <div key={byte.id} className="pt-8">
          <a href={byte.url} className="block text-xl font-semibold hover-underline text-base-content">
            {byte.title}
          </a>
          <p className="mt-1 text-sm text-primary">
            <time dateTime={byte.created}>{created.toFormat("d LLL yyyy")}</time>
          </p>
          <p className="mt-2 text-base text-base-content">{byte.summary}</p>
        </div>
      )
    })}
  </>
)

export default BytesList
