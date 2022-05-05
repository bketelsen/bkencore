import { DateTime } from "luxon"
import { FC } from "react"
import { bytes } from "../client/client"

const BytesList: FC<{ bytes: bytes.Byte[] }> = ({ bytes }) => (
  <>
    {bytes.map((byte) => {
      const created = DateTime.fromISO(byte.Created)
      return (
        <div key={byte.ID} className="pt-8">
          <a href={byte.URL} className="block text-xl font-semibold hover-underline text-neutral-900">
            {byte.Title}
          </a>
          <p className="mt-1 text-sm text-neutral-500">
            <time dateTime={byte.Created}>{created.toFormat("d LLL yyyy")}</time>
          </p>
          <p className="mt-2 text-base text-gray-600">{byte.Summary}</p>
        </div>
      )
    })}
  </>
)

export default BytesList
