import Client from "./client";

const env = process.env.NEXT_PUBLIC_ENCORE_ENV ?? "local"

export const DefaultClient = new Client(env)
