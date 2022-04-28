import { FC, useState } from "react"
import { DefaultClient } from "../client/default"

interface LoadingResult {
  result: "loading"
}

interface SuccessResult {
  result: "success"
}

interface ErrorResult {
  result: "error"
  err: Error
}

type Result = LoadingResult | SuccessResult | ErrorResult

const NewsletterSignup: FC = () => {
  const [result, setResult] = useState<Result | null>(null)
  const [email, setEmail] = useState("")

  const subscribe = async () => {
    if (email === "") { return }

    setResult({result: "loading"})
    try {
      await DefaultClient.email.Subscribe({Email: email})
      setResult({result: "success"})
    } catch(err) {
      setResult({result: "error", err: err as Error})
    }
  }

  return (
    <div>
      <div className="max-w-7xl mx-auto px-4 pt-12 sm:px-6 lg:pt-16 lg:px-8">
        <div className="px-6 py-6 bg-purple-700 rounded-lg md:py-12 md:px-12 lg:py-16 lg:px-16 xl:flex xl:items-center">
          <div className="xl:w-0 xl:flex-1">
            <h2 className="text-2xl font-extrabold tracking-tight text-white sm:text-3xl">
              {result?.result === "success" ? "Thanks for subscribing!" : "Want updates in your Inbox?"}
            </h2>
            <p className="mt-3 max-w-3xl text-lg leading-6 text-indigo-200">
              {result?.result === "success" ? "I'll keep you up to date with new posts." : "Sign up for my newsletter to stay up to date. No spam, I promise!"}
            </p>
          </div>
          <div className={`mt-8 sm:w-full sm:max-w-md xl:mt-0 xl:ml-8 ${result?.result === "success" ? "hidden xl:block xl:invisible" : ""}`}>
            <form className="sm:flex">
              <label htmlFor="email-address" className="sr-only">
                Email address
              </label>
              <input
                name="email-address"
                type="email"
                autoComplete="email"
                required
                className="w-full border-white px-5 py-3 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-700 focus:ring-white rounded-md"
                placeholder="Enter your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <button
                type="submit"
                className="mt-3 lg:w-36 flex items-center justify-center px-5 py-3 border border-transparent shadow text-base font-medium rounded-md text-white bg-purple-500 hover:bg-purple-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-indigo-700 focus:ring-white sm:mt-0 sm:ml-3 sm:w-auto sm:flex-shrink-0"
                disabled={result?.result === "loading"}
                onClick={subscribe}
              >
                {result?.result === "loading" &&
                  <svg className="flex-none animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
                  </svg>
                }
                Subscribe
              </button>
            </form>
            <p className="mt-3 text-sm text-indigo-200">
              I care about the protection of your data. Read my{' '}
              <a href="#" className="text-white font-medium underline">
                Privacy Policy.
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}

export default NewsletterSignup
