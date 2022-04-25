import type { NextPage } from 'next'
import Head from 'next/head'

const About: NextPage = () => {
  return (
    <div>
      <Head>
        <title>About Me | Brian Ketelsen</title>
      </Head>

      <p className="text-base lg:text-lg tracking-tight text-neutral-400">More than you probably need to know</p>
      <h1 className="text-4xl font-extrabold text-neutral-900">Oversharing</h1>
      <div className="mt-6 prose prose-indigo text-gray-500 max-w-prose">

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
  )
}

export default About
