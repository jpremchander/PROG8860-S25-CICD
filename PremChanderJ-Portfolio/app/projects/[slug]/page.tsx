import Link from "next/link"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { ArrowLeft, ExternalLink } from "lucide-react"

// This is a mock data function that would be replaced with real data in a production app
function getProjectData(slug: string) {
  // In a real application, you would fetch project data based on the slug
  const projects = {
    "aws-cost-optimization": {
      title: "AWS Cost Optimization",
      description:
        "Implemented comprehensive cost optimization strategies across multiple AWS accounts, reducing monthly spend by 30%.",
      client: "Cogniv Technologies",
      duration: "6 months",
      role: "Platform Engineer",
      tags: ["AWS", "Cost Management", "CloudWatch", "Resource Optimization"],
      overview:
        "As cloud costs continued to rise, our organization needed a systematic approach to identify waste and optimize spending across multiple AWS accounts without compromising performance or reliability.",
      challenge:
        "The main challenges included identifying unused or underutilized resources across multiple accounts, establishing governance for resource provisioning, and creating automated processes for ongoing optimization.",
      process: [
        {
          title: "Assessment",
          description:
            "Conducted a comprehensive audit of all AWS accounts, analyzing usage patterns and identifying potential savings opportunities.",
        },
        {
          title: "Resource Tagging Strategy",
          description:
            "Implemented a consistent tagging strategy to track resource ownership, purpose, and cost allocation.",
        },
        {
          title: "Automated Reporting",
          description:
            "Developed custom CloudWatch dashboards and automated reports to track spending trends and anomalies.",
        },
        {
          title: "Right-sizing Initiative",
          description:
            "Analyzed EC2 instance utilization and implemented right-sizing recommendations, converting to appropriate instance types.",
        },
        {
          title: "Reserved Instance Strategy",
          description: "Implemented a Reserved Instance purchasing strategy based on consistent usage patterns.",
        },
      ],
      outcome:
        "The project resulted in a 30% reduction in monthly AWS costs without impacting performance. We established automated monitoring for cost anomalies, implemented governance policies for new resource provisioning, and created a culture of cost awareness across engineering teams.",
    },
    "ci-cd-pipeline": {
      title: "CI/CD Pipeline Implementation",
      description: "Designed and implemented automated CI/CD pipelines using Jenkins, reducing deployment time by 70%.",
      client: "Maveric Systems",
      duration: "4 months",
      role: "DevOps Engineer",
      tags: ["Jenkins", "CI/CD", "Automation", "DevOps"],
      overview:
        "The development team was struggling with manual deployment processes that were error-prone and time-consuming. We needed to implement a robust CI/CD pipeline to automate the build, test, and deployment processes.",
      challenge:
        "The main challenges included integrating with existing systems, ensuring security throughout the pipeline, and minimizing disruption during the transition from manual to automated deployments.",
      process: [
        {
          title: "Requirements Analysis",
          description:
            "Worked with development and operations teams to understand the current workflow and identify requirements for the CI/CD pipeline.",
        },
        {
          title: "Tool Selection",
          description:
            "Evaluated various CI/CD tools and selected Jenkins based on our specific requirements and existing infrastructure.",
        },
        {
          title: "Pipeline Design",
          description:
            "Designed a pipeline architecture with stages for code checkout, build, test, security scanning, and deployment.",
        },
        {
          title: "Implementation",
          description:
            "Set up Jenkins servers, configured pipelines as code, and integrated with source control, testing frameworks, and deployment targets.",
        },
        {
          title: "Training and Documentation",
          description:
            "Created comprehensive documentation and conducted training sessions for development and operations teams.",
        },
      ],
      outcome:
        "The implementation reduced deployment time by 70% and virtually eliminated deployment-related errors. Developer productivity increased as they could focus on coding rather than deployment tasks. The system provided complete visibility into the deployment process and automated quality checks.",
    },
  }

  return projects[slug as keyof typeof projects] || projects["aws-cost-optimization"]
}

export default function ProjectPage({ params }: { params: { slug: string } }) {
  const projectData = getProjectData(params.slug)

  return (
    <div className="flex min-h-screen flex-col">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container flex h-16 items-center justify-between">
          <Link href="/" className="flex items-center gap-2 text-lg font-bold">
            <span>Prem Chander J</span>
          </Link>
          <nav className="hidden md:flex gap-6">
            <Link href="/#work" className="text-sm font-medium transition-colors hover:text-primary">
              Work
            </Link>
            <Link href="/#about" className="text-sm font-medium transition-colors hover:text-primary">
              About
            </Link>
            <Link href="/#skills" className="text-sm font-medium transition-colors hover:text-primary">
              Skills
            </Link>
            <Link href="/#contact" className="text-sm font-medium transition-colors hover:text-primary">
              Contact
            </Link>
          </nav>
          <Button size="sm" className="hidden md:flex">
            Resume
          </Button>
          <Button variant="ghost" size="icon" className="md:hidden">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="h-6 w-6"
            >
              <line x1="4" x2="20" y1="12" y2="12" />
              <line x1="4" x2="20" y1="6" y2="6" />
              <line x1="4" x2="20" y1="18" y2="18" />
            </svg>
          </Button>
        </div>
      </header>
      <main className="flex-1">
        <div className="container py-12 md:py-24">
          <Button variant="ghost" size="sm" className="mb-8" asChild>
            <Link href="/projects">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to Projects
            </Link>
          </Button>

          <div className="relative aspect-[21/9] w-full overflow-hidden rounded-lg">
            <Image
              src="/placeholder.svg?height=900&width=1900"
              alt={projectData.title}
              width={1900}
              height={900}
              className="object-cover"
              priority
            />
          </div>

          <div className="mx-auto max-w-3xl py-12">
            <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">{projectData.title}</h1>
            <p className="mt-4 text-xl text-muted-foreground">{projectData.description}</p>

            <div className="mt-8 grid gap-4 sm:grid-cols-2 md:grid-cols-3">
              <div>
                <h3 className="font-medium">Client</h3>
                <p className="text-muted-foreground">{projectData.client}</p>
              </div>
              <div>
                <h3 className="font-medium">Duration</h3>
                <p className="text-muted-foreground">{projectData.duration}</p>
              </div>
              <div>
                <h3 className="font-medium">Role</h3>
                <p className="text-muted-foreground">{projectData.role}</p>
              </div>
            </div>

            <div className="mt-4 flex flex-wrap gap-2">
              {projectData.tags.map((tag: string, index: number) => (
                <div key={index} className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">
                  {tag}
                </div>
              ))}
            </div>

            <div className="mt-12">
              <h2 className="text-2xl font-bold">Overview</h2>
              <p className="mt-4 text-muted-foreground">{projectData.overview}</p>
            </div>

            <div className="mt-12">
              <h2 className="text-2xl font-bold">Challenge</h2>
              <p className="mt-4 text-muted-foreground">{projectData.challenge}</p>
            </div>

            <div className="mt-12">
              <h2 className="text-2xl font-bold">Process</h2>
              <div className="mt-6 grid gap-8">
                {projectData.process.map((step: { title: string; description: string }, index: number) => (
                  <div key={index} className="flex gap-4">
                    <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary text-primary-foreground">
                      {index + 1}
                    </div>
                    <div>
                      <h3 className="font-bold">{step.title}</h3>
                      <p className="mt-2 text-muted-foreground">{step.description}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="mt-12 grid gap-6 md:grid-cols-2">
              <div className="overflow-hidden rounded-lg">
                <Image
                  src="/placeholder.svg?height=600&width=800"
                  alt="Project diagram"
                  width={800}
                  height={600}
                  className="object-cover"
                />
              </div>
              <div className="overflow-hidden rounded-lg">
                <Image
                  src="/placeholder.svg?height=600&width=800"
                  alt="Project implementation"
                  width={800}
                  height={600}
                  className="object-cover"
                />
              </div>
            </div>

            <div className="mt-12">
              <h2 className="text-2xl font-bold">Outcome</h2>
              <p className="mt-4 text-muted-foreground">{projectData.outcome}</p>
            </div>

            <div className="mt-12 flex justify-center">
              <Button asChild>
                <Link href="/#contact">
                  <ExternalLink className="mr-2 h-4 w-4" />
                  Contact Me for Similar Projects
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
      <footer className="border-t py-6 md:py-0">
        <div className="container flex flex-col items-center justify-between gap-4 md:h-16 md:flex-row">
          <p className="text-sm text-muted-foreground">
            © {new Date().getFullYear()} Prem Chander J. All rights reserved.
          </p>
          <div className="flex gap-4">
            <Link href="#" className="text-sm text-muted-foreground hover:text-primary">
              Privacy Policy
            </Link>
            <Link href="#" className="text-sm text-muted-foreground hover:text-primary">
              Terms of Service
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
