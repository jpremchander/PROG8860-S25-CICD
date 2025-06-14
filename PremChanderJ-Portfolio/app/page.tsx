"use client"

import type React from "react"

import { useState, useEffect } from "react"
import Image from "next/image"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Textarea } from "@/components/ui/textarea"
import { useToast } from "@/hooks/use-toast"
import { Download, Send, Linkedin, Mail, Phone, MapPin, Menu, X, Github } from "lucide-react"

export default function Home() {
  const { toast } = useToast()
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [message, setMessage] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false)
  const [activeSection, setActiveSection] = useState("home")

  useEffect(() => {
    const handleScroll = () => {
      const sections = ["home", "about", "work", "contact"]

      for (const section of sections) {
        const element = document.getElementById(section)
        if (element) {
          const rect = element.getBoundingClientRect()
          if (rect.top <= 100 && rect.bottom >= 100) {
            setActiveSection(section)
            break
          }
        }
      }
    }

    window.addEventListener("scroll", handleScroll)
    return () => window.removeEventListener("scroll", handleScroll)
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      // In a real implementation, you would send this data to your backend
      // For now, we'll simulate a successful submission
      await new Promise((resolve) => setTimeout(resolve, 1000))

      toast({
        title: "Message sent!",
        description: "Thank you for reaching out. I'll get back to you soon.",
      })

      // Reset form
      setName("")
      setEmail("")
      setMessage("")
    } catch {
      toast({
        title: "Something went wrong",
        description: "Your message couldn't be sent. Please try again.",
        variant: "destructive",
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  const scrollToSection = (sectionId: string) => {
    const element = document.getElementById(sectionId)
    if (element) {
      element.scrollIntoView({ behavior: "smooth" })
    }
    setMobileMenuOpen(false)
  }

  return (
    <div className="min-h-screen flex flex-col dark">
      {/* Header */}
      <header className="sticky top-0 z-50 backdrop-blur-md bg-background/80 border-b border-border">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-xl font-bold gradient-text">
            Prem Chander J
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            <button
              onClick={() => scrollToSection("home")}
              className={`text-sm font-medium transition-colors hover:text-primary ${activeSection === "home" ? "text-primary" : "text-muted-foreground"}`}
            >
              Home
            </button>
            <button
              onClick={() => scrollToSection("about")}
              className={`text-sm font-medium transition-colors hover:text-primary ${activeSection === "about" ? "text-primary" : "text-muted-foreground"}`}
            >
              About
            </button>
            <button
              onClick={() => scrollToSection("work")}
              className={`text-sm font-medium transition-colors hover:text-primary ${activeSection === "work" ? "text-primary" : "text-muted-foreground"}`}
            >
              Work
            </button>
            <button
              onClick={() => scrollToSection("contact")}
              className={`text-sm font-medium transition-colors hover:text-primary ${activeSection === "contact" ? "text-primary" : "text-muted-foreground"}`}
            >
              Contact
            </button>
            <Button asChild size="sm">
              <a href="/resume.pdf" download>
                <Download className="mr-2 h-4 w-4" />
                Resume
              </a>
            </Button>
          </nav>

          {/* Mobile Menu Button */}
          <Button variant="ghost" size="icon" className="md:hidden" onClick={() => setMobileMenuOpen(!mobileMenuOpen)}>
            {mobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
          </Button>
        </div>

        {/* Mobile Navigation */}
        {mobileMenuOpen && (
          <div className="md:hidden absolute w-full bg-background border-b border-border py-4 px-4 flex flex-col space-y-4">
            <button onClick={() => scrollToSection("home")} className="text-sm font-medium py-2 hover:text-primary">
              Home
            </button>
            <button onClick={() => scrollToSection("about")} className="text-sm font-medium py-2 hover:text-primary">
              About
            </button>
            <button onClick={() => scrollToSection("work")} className="text-sm font-medium py-2 hover:text-primary">
              Work
            </button>
            <button onClick={() => scrollToSection("contact")} className="text-sm font-medium py-2 hover:text-primary">
              Contact
            </button>
            <Button asChild size="sm" className="w-full">
              <a href="/resume.pdf" download>
                <Download className="mr-2 h-4 w-4" />
                Resume
              </a>
            </Button>
          </div>
        )}
      </header>

      <main className="flex-1">
        {/* Hero Section */}
        <section id="home" className="py-20 md:py-32">
          <div className="container mx-auto px-4">
            <div className="flex flex-col md:flex-row items-center gap-8 md:gap-16">
              <div className="flex-1 space-y-6">
                <h1 className="text-4xl md:text-6xl font-bold leading-tight">
                  Hi, I&apos;m <span className="gradient-text">Prem Chander J</span>
                </h1>
                <h2 className="text-2xl md:text-3xl font-medium text-muted-foreground">Platform Engineer</h2>
                <p className="text-lg text-muted-foreground max-w-xl">
                  Results-driven Platform Engineer with extensive expertise in cloud architecture, DevOps methodologies,
                  and infrastructure optimization. Specializing in AWS, Azure, and Huawei Cloud environments with a
                  proven track record of implementing cost-effective, scalable solutions.
                </p>
                <div className="flex flex-wrap gap-4 pt-4">
                  <Button onClick={() => scrollToSection("about")}>Learn More About Me</Button>
                  <Button variant="outline" onClick={() => scrollToSection("contact")}>
                    Get In Touch
                  </Button>
                </div>
              </div>
              <div className="relative w-64 h-64 md:w-80 md:h-80 rounded-full overflow-hidden border-4 border-primary/20">
                <Image src="/images/profile.png" alt="Prem Chander J" fill className="object-cover" priority />
              </div>
            </div>
          </div>
        </section>

        {/* About Section */}
        <section id="about" className="py-20 bg-secondary/20">
          <div className="container mx-auto px-4">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl font-bold mb-8 text-center">About Me</h2>

              <div className="space-y-6">
                <p className="text-lg">
                  I&apos;m a Platform Engineer with extensive experience in cloud infrastructure, DevOps practices, and
                  monitoring solutions. My expertise spans across AWS, Azure, and Huawei Cloud platforms, with a focus
                  on automation, optimization, and security.
                </p>

                <div className="grid md:grid-cols-2 gap-6">
                  <div className="gradient-border p-6 bg-card">
                    <h3 className="text-xl font-semibold mb-4">Experience</h3>
                    <div className="space-y-4">
                      <div>
                        <h4 className="font-medium">Platform Engineer</h4>
                        <p className="text-sm text-muted-foreground">Cogniv Technologies | 2017 - Present</p>
                        <ul className="list-disc list-inside text-sm mt-2 text-muted-foreground">
                          <li>Analyzing and forecasting cloud costs for business planning</li>
                          <li>Monitoring cloud resource usage and optimizing for cost efficiency</li>
                          <li>Collaborating with cross-functional teams to optimize infrastructure</li>
                        </ul>
                      </div>
                      <div>
                        <h4 className="font-medium">DevOps Engineer</h4>
                        <p className="text-sm text-muted-foreground">Maveric Systems Pvt Ltd</p>
                        <ul className="list-disc list-inside text-sm mt-2 text-muted-foreground">
                          <li>Managed AWS infrastructure including EC2, S3, and RDS</li>
                          <li>Configured load balancers and auto-scaling for high availability</li>
                          <li>Led cloud cost optimization initiatives</li>
                        </ul>
                      </div>
                    </div>
                  </div>

                  <div className="gradient-border p-6 bg-card">
                    <h3 className="text-xl font-semibold mb-4">Skills & Certifications</h3>
                    <div className="space-y-4">
                      <div>
                        <h4 className="font-medium">Cloud Platforms</h4>
                        <p className="text-sm text-muted-foreground mt-1">AWS, Azure, Huawei Cloud</p>
                      </div>
                      <div>
                        <h4 className="font-medium">DevOps Tools</h4>
                        <p className="text-sm text-muted-foreground mt-1">
                          Jenkins, Spinnaker, Ansible, Git, Bitbucket
                        </p>
                      </div>
                      <div>
                        <h4 className="font-medium">Certifications</h4>
                        <ul className="list-disc list-inside text-sm text-muted-foreground">
                          <li>RHCSA/RHCE (Certificate #: 130-159-836)</li>
                          <li>AWS Certified Solution Architect – Associate</li>
                          <li>AWS Certified Solution Professional & DevOps (Pursuing)</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="flex justify-center pt-6">
                  <Button asChild>
                    <a href="/resume.pdf" download>
                      <Download className="mr-2 h-4 w-4" />
                      Download Full Resume
                    </a>
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Work Section */}
        <section id="work" className="py-20">
          <div className="container mx-auto px-4">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl font-bold mb-4 text-center">My Work Experience</h2>
              <p className="text-lg text-muted-foreground mb-12 text-center">
                A summary of my professional journey and key projects
              </p>

              <div className="space-y-12">
                {/* Cogniv */}
                <div className="gradient-border p-6 bg-card">
                  <div className="flex flex-col md:flex-row justify-between mb-4">
                    <h3 className="text-xl font-semibold">Cogniv Technologies</h3>
                    <p className="text-sm text-muted-foreground">2017 - Present</p>
                  </div>
                  <p className="font-medium mb-2">Platform Engineer</p>
                  <ul className="list-disc list-inside text-muted-foreground space-y-2">
                    <li>Analyzing and forecasting cloud costs to help businesses plan and budget their IT expenses</li>
                    <li>Monitoring cloud resource usage, identifying unused or underutilized resources</li>
                    <li>Collaborating with different teams to optimize cloud infrastructure and reduce costs</li>
                    <li>
                      Monitoring and analyzing cloud-related expenses to spot trends and identify savings opportunities
                    </li>
                  </ul>
                </div>

                {/* Maveric */}
                <div className="gradient-border p-6 bg-card">
                  <div className="flex flex-col md:flex-row justify-between mb-4">
                    <h3 className="text-xl font-semibold">Maveric Systems Pvt Ltd</h3>
                    <p className="text-sm text-muted-foreground">Previous Role</p>
                  </div>
                  <p className="font-medium mb-2">DevOps Engineer</p>
                  <ul className="list-disc list-inside text-muted-foreground space-y-2">
                    <li>Actively managed, improvised, and monitored cloud infrastructure on AWS, EC2, S3, and RDS</li>
                    <li>Created and configured elastic load balancers and auto scaling groups</li>
                    <li>Maintained the DevOps GitHub Repo and its branches for scripts and tasks</li>
                    <li>Effectively managed Cloud Cost Optimization across AWS Accounts</li>
                  </ul>
                </div>

                {/* Chimera */}
                <div className="gradient-border p-6 bg-card">
                  <div className="flex flex-col md:flex-row justify-between mb-4">
                    <h3 className="text-xl font-semibold">Chimera Technologies Pvt Ltd</h3>
                    <p className="text-sm text-muted-foreground">Previous Role</p>
                  </div>
                  <p className="font-medium mb-2">DevOps Engineer</p>
                  <ul className="list-disc list-inside text-muted-foreground space-y-2">
                    <li>Migrated AWS Environment to Huawei Cloud</li>
                    <li>Worked on Architectural Designs for infrastructure using draw.io tools</li>
                    <li>Managed documentation and sprint tasks using Confluence and Jira</li>
                    <li>Created CloudFormation Templates for VPC and infrastructure resources</li>
                  </ul>
                </div>

                {/* Match Move */}
                <div className="gradient-border p-6 bg-card">
                  <div className="flex flex-col md:flex-row justify-between mb-4">
                    <h3 className="text-xl font-semibold">Match Move India Pvt Ltd (FIN-TECH)</h3>
                    <p className="text-sm text-muted-foreground">Aug 2019 – July 2021</p>
                  </div>
                  <p className="font-medium mb-2">DevOps Engineer / SRE</p>
                  <ul className="list-disc list-inside text-muted-foreground space-y-2">
                    <li>Established processes for maintenance, monitoring, and security patching</li>
                    <li>Worked for PCI DSS Audit Certificate and secured infrastructure and APIs</li>
                    <li>Improved performance and security for high-traffic websites</li>
                    <li>Managed a team of SRE Engineers for SLAs, RCAs, Monitoring, and Reporting</li>
                  </ul>
                </div>
              </div>

              <div className="mt-12 grid md:grid-cols-3 gap-6">
                <div className="p-6 rounded-lg bg-card border border-border">
                  <h3 className="font-semibold mb-2">Cloud Infrastructure</h3>
                  <p className="text-sm text-muted-foreground">
                    Designing and managing scalable, secure cloud environments on AWS, Azure, and Huawei Cloud.
                  </p>
                </div>
                <div className="p-6 rounded-lg bg-card border border-border">
                  <h3 className="font-semibold mb-2">DevOps Automation</h3>
                  <p className="text-sm text-muted-foreground">
                    Implementing CI/CD pipelines and infrastructure as code for streamlined deployments.
                  </p>
                </div>
                <div className="p-6 rounded-lg bg-card border border-border">
                  <h3 className="font-semibold mb-2">Cost Optimization</h3>
                  <p className="text-sm text-muted-foreground">
                    Analyzing and optimizing cloud resources to maximize efficiency and minimize expenses.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Skills Section */}
        <section className="py-20 bg-secondary/20">
          <div className="container mx-auto px-4">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl font-bold mb-8 text-center">Technical Skills</h2>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="gradient-border p-6 bg-card">
                  <h3 className="text-xl font-semibold mb-4">Cloud Platforms & Infrastructure</h3>
                  <ul className="list-disc list-inside text-muted-foreground space-y-1">
                    <li>AWS (EC2, S3, RDS, VPC, Auto Scaling, Load Balancer)</li>
                    <li>Azure</li>
                    <li>Huawei Cloud</li>
                    <li>CloudWatch, Inspector, Systems Manager</li>
                    <li>Route 53, API Gateway, IAM, WAF, Guard Duty</li>
                    <li>Lambda, CloudFront</li>
                  </ul>
                </div>

                <div className="gradient-border p-6 bg-card">
                  <h3 className="text-xl font-semibold mb-4">DevOps & Automation</h3>
                  <ul className="list-disc list-inside text-muted-foreground space-y-1">
                    <li>Jenkins, Spinnaker, Ansible</li>
                    <li>Python (boto3)</li>
                    <li>Shell scripting</li>
                    <li>Bitbucket and Git</li>
                    <li>JIRA, Confluence</li>
                  </ul>
                </div>

                <div className="gradient-border p-6 bg-card">
                  <h3 className="text-xl font-semibold mb-4">Monitoring Tools</h3>
                  <ul className="list-disc list-inside text-muted-foreground space-y-1">
                    <li>CloudWatch, Prometheus, New Relic</li>
                    <li>Gray log, ELK</li>
                    <li>Pingdom, Sensu</li>
                    <li>Nagios, Cacti, Zabbix</li>
                    <li>Splunk</li>
                  </ul>
                </div>

                <div className="gradient-border p-6 bg-card">
                  <h3 className="text-xl font-semibold mb-4">Operating Systems</h3>
                  <ul className="list-disc list-inside text-muted-foreground space-y-1">
                    <li>Ubuntu 12.04/14.04/16.04/18.04</li>
                    <li>CentOS</li>
                    <li>RedHat Linux</li>
                    <li>Amazon Linux</li>
                    <li>Windows 2008 R2/XP/Win7/8/10</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Contact Section */}
        <section id="contact" className="py-20">
          <div className="container mx-auto px-4">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-3xl font-bold mb-8 text-center">Get In Touch</h2>

              <div className="grid md:grid-cols-2 gap-8">
                <div>
                  <h3 className="text-xl font-semibold mb-6">Contact Information</h3>
                  <div className="space-y-4">
                    <div className="flex items-center">
                      <Mail className="h-5 w-5 mr-3 text-primary" />
                      <a href="mailto:premchander.j.pc@gmail.com" className="text-muted-foreground hover:text-primary">
                        premchander.j.pc@gmail.com
                      </a>
                    </div>
                    <div className="flex items-center">
                      <Phone className="h-5 w-5 mr-3 text-primary" />
                      <a href="tel:+919600503030" className="text-muted-foreground hover:text-primary">
                        +91 9600503030
                      </a>
                    </div>
                    <div className="flex items-center">
                      <MapPin className="h-5 w-5 mr-3 text-primary" />
                      <span className="text-muted-foreground">India</span>
                    </div>
                  </div>

                  <h3 className="text-xl font-semibold mt-8 mb-6">Connect With Me</h3>
                  <div className="flex space-x-4">
                    <Button asChild variant="outline" size="icon">
                      <a
                        href="https://www.linkedin.com/in/premchander-j/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="LinkedIn Profile"
                      >
                        <Linkedin className="h-5 w-5" />
                      </a>
                    </Button>
                    <Button asChild variant="outline" size="icon">
                      <a
                        href="https://github.com/"
                        target="_blank"
                        rel="noopener noreferrer"
                        aria-label="GitHub Profile"
                      >
                        <Github className="h-5 w-5" />
                      </a>
                    </Button>
                    <Button asChild variant="outline" size="icon">
                      <a href="mailto:premchander.j.pc@gmail.com" aria-label="Email">
                        <Mail className="h-5 w-5" />
                      </a>
                    </Button>
                  </div>
                </div>

                <div className="gradient-border p-6 bg-card">
                  <h3 className="text-xl font-semibold mb-6">Send Me a Message</h3>
                  <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                      <label htmlFor="name" className="block text-sm font-medium mb-1">
                        Name
                      </label>
                      <Input
                        id="name"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="Your name"
                        required
                      />
                    </div>
                    <div>
                      <label htmlFor="email" className="block text-sm font-medium mb-1">
                        Email
                      </label>
                      <Input
                        id="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Your email"
                        required
                      />
                    </div>
                    <div>
                      <label htmlFor="message" className="block text-sm font-medium mb-1">
                        Message
                      </label>
                      <Textarea
                        id="message"
                        value={message}
                        onChange={(e) => setMessage(e.target.value)}
                        placeholder="Your message"
                        rows={4}
                        required
                      />
                    </div>
                    <Button type="submit" className="w-full" disabled={isSubmitting}>
                      {isSubmitting ? (
                        <span className="flex items-center">
                          <svg
                            className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                          >
                            <circle
                              className="opacity-25"
                              cx="12"
                              cy="12"
                              r="10"
                              stroke="currentColor"
                              strokeWidth="4"
                            ></circle>
                            <path
                              className="opacity-75"
                              fill="currentColor"
                              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            ></path>
                          </svg>
                          Sending...
                        </span>
                      ) : (
                        <span className="flex items-center">
                          <Send className="mr-2 h-4 w-4" />
                          Send Message
                        </span>
                      )}
                    </Button>
                  </form>
                </div>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="py-6 border-t border-border">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row justify-between items-center">
            <p className="text-sm text-muted-foreground">
              © {new Date().getFullYear()} Prem Chander J. All rights reserved.
            </p>
            <div className="flex space-x-6 mt-4 md:mt-0">
              <a
                href="https://www.linkedin.com/in/premchander-j/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary"
              >
                LinkedIn
              </a>
              <a
                href="https://github.com/"
                target="_blank"
                rel="noopener noreferrer"
                className="text-muted-foreground hover:text-primary"
              >
                GitHub
              </a>
              <a href="mailto:premchander.j.pc@gmail.com" className="text-muted-foreground hover:text-primary">
                Email
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
