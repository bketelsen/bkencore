import type { NextPage } from 'next'
import Head from 'next/head'

const Home: NextPage = () => {
  return (
    <div>
      <Head>
        <title>About Me | Brian Ketelsen</title>
      </Head>

      <div className="relative py-16 bg-white overflow-hidden">
        <div className="hidden lg:block lg:absolute lg:inset-y-0 lg:h-full lg:w-full">
          <div className="relative h-full text-lg max-w-prose mx-auto" aria-hidden="true">
            <svg
              className="absolute top-12 left-full transform translate-x-32"
              width={404}
              height={384}
              fill="none"
              viewBox="0 0 404 384"
            >
              <defs>
                <pattern
                  id="74b3fd99-0a6f-4271-bef2-e80eeafdf357"
                  x={0}
                  y={0}
                  width={20}
                  height={20}
                  patternUnits="userSpaceOnUse"
                >
                  <rect x={0} y={0} width={4} height={4} className="text-gray-200" fill="currentColor" />
                </pattern>
              </defs>
              <rect width={404} height={384} fill="url(#74b3fd99-0a6f-4271-bef2-e80eeafdf357)" />
            </svg>
            <svg
              className="absolute top-1/2 right-full transform -translate-y-1/2 -translate-x-32"
              width={404}
              height={384}
              fill="none"
              viewBox="0 0 404 384"
            >
              <defs>
                <pattern
                  id="f210dbf6-a58d-4871-961e-36d5016a0f49"
                  x={0}
                  y={0}
                  width={20}
                  height={20}
                  patternUnits="userSpaceOnUse"
                >
                  <rect x={0} y={0} width={4} height={4} className="text-gray-200" fill="currentColor" />
                </pattern>
              </defs>
              <rect width={404} height={384} fill="url(#f210dbf6-a58d-4871-961e-36d5016a0f49)" />
            </svg>

          </div>
        </div>
        <div className="relative px-4 sm:px-6 lg:px-8">
          <div className="text-lg max-w-prose mx-auto">
            <h1>
              <span className="mt-2 block text-3xl text-center leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl">
                Oversharing
              </span>
            </h1>
            <p className="mt-8 text-xl text-gray-500 leading-8">
              More than you probably need to know
            </p>
          </div>
          <div className="mt-6 prose prose-indigo prose-lg text-gray-500 mx-auto">

            <p>
              Howdy! Thanks for stopping by. I&apos;m Brian Ketelsen, and I&apos;ve been doing technology things
              since acoustic modem couplers were a thing. I love Open Source and exploring different
              programming languages. Some highlights about me and my activities:
            </p>

            <ul>
              <li>Co-founded GopherCon, the largest conference for Go developers</li>
              <li>Co-authored Go In Action for Manning Publishing</li>
              <li>Along with Erik St. Martin, I wrote Skynet and SkyDNS, which was eventually morphed into the DNS service discovery that powers Kubernetes.</li>
              <li>Contributor to virtual kubelet and Krustlet</li>
            </ul>

            <p>
              I&apos;ve worked in a lot of interesting industries from Consumer Credit to Healthcare,
              in jobs ranging from DBA to Chief Information Officer. I&apos;m currently loving life
              at Microsoft as a Cloud Developer Advocate.
            </p>
            <p>
              I love to teach because I&apos;m always learning new things. Cloud Advocacy also
              gives me the freedom to experiment with new and interesting things every day.
            </p>
            <p>
              My most exciting moment in Open Source came when I got a letter of thanks from
              the JPL at NASA for a library I wrote. They used it on one of the Rover missions.
            </p>

          </div>
        </div>
      </div>

    </div>
  )
}

export default Home
