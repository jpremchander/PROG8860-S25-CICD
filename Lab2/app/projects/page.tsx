import Link from "next/link"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { ArrowLeft } from "lucide-react"

export default function ProjectsPage() {
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
        <section className="container py-12 md:py-24">
          <div className="flex flex-col gap-4">
            <Button variant="ghost" size="sm" className="w-fit" asChild>
              <Link href="/">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Back to Home
              </Link>
            </Button>
            <h1 className="text-3xl font-bold tracking-tighter sm:text-4xl md:text-5xl">Technical Projects</h1>
            <p className="max-w-[85%] text-muted-foreground sm:text-lg">
              A collection of my work across cloud infrastructure, DevOps, and automation.
            </p>
          </div>
          <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3 py-12">
            {/* Project 1 */}
            <div className="group relative overflow-hidden rounded-lg border">
              <Link href="/projects/aws-cost-optimization" className="absolute inset-0 z-10">
                <span className="sr-only">View Project</span>
              </Link>
              <div className="relative aspect-[16/10] overflow-hidden">
                <Image
                  src="/placeholder.svg?height=450&width=720"
                  alt="AWS Cost Optimization"
                  width={720}
                  height={450}
                  className="object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="p-6">
                <h3 className="text-xl font-bold">AWS Cost Optimization</h3>
                <p className="mt-2 text-muted-foreground">
                  Implemented comprehensive cost optimization strategies across multiple AWS accounts, reducing monthly
                  spend by 30%.
                </p>
                <div className="mt-4 flex flex-wrap gap-2">
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">AWS</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Cost Management</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">CloudWatch</div>
                </div>
              </div>
            </div>

            {/* Project 2 */}
            <div className="group relative overflow-hidden rounded-lg border">
              <Link href="/projects/ci-cd-pipeline" className="absolute inset-0 z-10">
                <span className="sr-only">View Project</span>
              </Link>
              <div className="relative aspect-[16/10] overflow-hidden">
                <Image
                  src="/placeholder.svg?height=450&width=720"
                  alt="CI/CD Pipeline Implementation"
                  width={720}
                  height={450}
                  className="object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="p-6">
                <h3 className="text-xl font-bold">CI/CD Pipeline Implementation</h3>
                <p className="mt-2 text-muted-foreground">
                  Designed and implemented automated CI/CD pipelines using Jenkins, reducing deployment time by 70%.
                </p>
                <div className="mt-4 flex flex-wrap gap-2">
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Jenkins</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">CI/CD</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Automation</div>
                </div>
              </div>
            </div>

            {/* Project 3 */}
            <div className="group relative overflow-hidden rounded-lg border">
              <Link href="/projects/cloud-migration" className="absolute inset-0 z-10">
                <span className="sr-only">View Project</span>
              </Link>
              <div className="relative aspect-[16/10] overflow-hidden">
                <Image
                  src="/placeholder.svg?height=450&width=720"
                  alt="Cloud Migration Project"
                  width={720}
                  height={450}
                  className="object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="p-6">
                <h3 className="text-xl font-bold">Cloud Migration Project</h3>
                <p className="mt-2 text-muted-foreground">
                  Successfully migrated on-premises infrastructure to AWS, improving scalability and reducing
                  operational costs.
                </p>
                <div className="mt-4 flex flex-wrap gap-2">
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Cloud Migration</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">AWS</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Infrastructure</div>
                </div>
              </div>
            </div>

            {/* Project 4 */}
            <div className="group relative overflow-hidden rounded-lg border">
              <Link href="/projects/monitoring-solution" className="absolute inset-0 z-10">
                <span className="sr-only">View Project</span>
              </Link>
              <div className="relative aspect-[16/10] overflow-hidden">
                <Image
                  src="/placeholder.svg?height=450&width=720"
                  alt="Enterprise Monitoring Solution"
                  width={720}
                  height={450}
                  className="object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="p-6">
                <h3 className="text-xl font-bold">Enterprise Monitoring Solution</h3>
                <p className="mt-2 text-muted-foreground">
                  Implemented comprehensive monitoring using Prometheus, ELK, and CloudWatch, improving incident
                  response time by 60%.
                </p>
                <div className="mt-4 flex flex-wrap gap-2">
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Monitoring</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Prometheus</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">ELK Stack</div>
                </div>
              </div>
            </div>

            {/* Project 5 */}
            <div className="group relative overflow-hidden rounded-lg border">
              <Link href="/projects/security-compliance" className="absolute inset-0 z-10">
                <span className="sr-only">View Project</span>
              </Link>
              <div className="relative aspect-[16/10] overflow-hidden">
                <Image
                  src="/placeholder.svg?height=450&width=720"
                  alt="Security & Compliance Implementation"
                  width={720}
                  height={450}
                  className="object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="p-6">
                <h3 className="text-xl font-bold">Security & Compliance Implementation</h3>
                <p className="mt-2 text-muted-foreground">
                  Led PCI DSS compliance initiative, implementing security best practices across cloud infrastructure.
                </p>
                <div className="mt-4 flex flex-wrap gap-2">
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Security</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">PCI DSS</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Compliance</div>
                </div>
              </div>
            </div>

            {/* Project 6 */}
            <div className="group relative overflow-hidden rounded-lg border">
              <Link href="/projects/infrastructure-automation" className="absolute inset-0 z-10">
                <span className="sr-only">View Project</span>
              </Link>
              <div className="relative aspect-[16/10] overflow-hidden">
                <Image
                  src="/placeholder.svg?height=450&width=720"
                  alt="Infrastructure Automation"
                  width={720}
                  height={450}
                  className="object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>
              <div className="p-6">
                <h3 className="text-xl font-bold">Infrastructure Automation</h3>
                <p className="mt-2 text-muted-foreground">
                  Developed CloudFormation templates and Ansible playbooks for automated infrastructure provisioning.
                </p>
                <div className="mt-4 flex flex-wrap gap-2">
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">CloudFormation</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">Ansible</div>
                  <div className="rounded-full bg-muted px-2.5 py-0.5 text-xs font-semibold">IaC</div>
                </div>
              </div>
            </div>
          </div>
        </section>
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
