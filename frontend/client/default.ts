import Client from "./client";

const env = process.env.NEXT_PUBLIC_ENCORE_ENV ?? "staging"

export const DefaultClient = new Client(env)
