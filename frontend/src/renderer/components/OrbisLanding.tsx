import { useState } from 'react';
import { Mail } from 'lucide-react';
import { Link } from '@tanstack/react-router';

const TwitterIcon = ({ className }: { className?: string }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}><path d="M22 4s-.7 2.1-2 3.4c1.6 10-9.4 17.3-18 11.6 2.2.1 4.4-.6 6-2C3 15.5.5 9.6 3 5c2.2 2.6 5.6 4.1 9 4-.9-4.2 4-6.6 7-3.8 1.1 0 3-1.2 3-1.2z"/></svg>
);

const GithubIcon = ({ className }: { className?: string }) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className}><path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/></svg>
);

const FaqItem = ({ q, a }: { q: string, a: string }) => {
  const [open, setOpen] = useState(false);
  return (
    <div className="border-b border-white/10 pb-6 group cursor-pointer" onClick={() => setOpen(!open)}>
      <div className="flex justify-between items-center mb-2">
        <h3 className={`font-mono text-[16px] sm:text-[20px] uppercase transition-colors pr-8 ${open ? 'text-cream' : 'text-neon group-hover:text-cream'}`}>{q}</h3>
        <span className={`transition-colors text-[24px] ${open ? 'text-cream' : 'text-neon group-hover:text-cream'}`}>{open ? '-' : '+'}</span>
      </div>
      {open && <p className="font-mono text-[14px] text-cream/70 mt-4 leading-relaxed max-w-[80%]">{a}</p>}
    </div>
  );
};

export function OrbisLanding() {
  const [activeTab, setActiveTab] = useState('homepage');

  return (
    <div className="h-screen w-full overflow-y-auto scroll-smooth bg-[#010828] text-cream font-mono relative overflow-x-hidden selection:bg-neon selection:text-[#010828]">
      
      {/* Texture Overlay */}
      <div 
        className="fixed inset-0 z-50 pointer-events-none bg-cover mix-blend-lighten opacity-60" 
        style={{ backgroundImage: "url('/texture.png')" }} 
      />

      <div className="w-full max-w-[1831px] mx-auto">
        {/* Header / Navbar (Moved outside section to prevent clipping) */}
        <header className="flex items-center justify-between py-8 z-[100] fixed top-4 sm:top-8 left-0 right-0 px-6 sm:px-12 lg:px-24 pointer-events-none">
          <div className="font-grotesk text-[40px] uppercase tracking-wider text-cream pointer-events-auto liquid-glass rounded-[28px] px-[40px] py-[16px] lg:px-[52px] lg:py-[24px] flex items-center justify-center leading-none">
            Drake
          </div>
          
          <nav className="hidden lg:flex items-center gap-2 liquid-glass rounded-[28px] px-[40px] py-[16px] pointer-events-auto">
            <a 
              href="#" 
              onClick={() => setActiveTab('homepage')}
              className={`font-grotesk text-[13px] uppercase px-6 py-3 rounded-full transition-all duration-500 ease-out ${activeTab === 'homepage' ? 'text-[#010828] bg-neon' : 'text-cream hover:text-neon bg-transparent'}`}
            >
              Homepage
            </a>
            <a 
              href="#agents" 
              onClick={() => setActiveTab('agents')}
              className={`font-grotesk text-[13px] uppercase px-6 py-3 rounded-full transition-all duration-500 ease-out ${activeTab === 'agents' ? 'text-[#010828] bg-neon' : 'text-cream hover:text-neon bg-transparent'}`}
            >
              Agents
            </a>
            <a 
              href="#faq" 
              onClick={() => setActiveTab('faq')}
              className={`font-grotesk text-[13px] uppercase px-6 py-3 rounded-full transition-all duration-500 ease-out ${activeTab === 'faq' ? 'text-[#010828] bg-neon' : 'text-cream hover:text-neon bg-transparent'}`}
            >
              FAQ
            </a>
            <div className="w-[1px] h-[24px] bg-white/20 mx-2" />
            <Link to="/app" className="font-grotesk text-[13px] uppercase text-neon hover:text-cream transition-colors px-4">
              Launch Workspace
            </Link>
          </nav>
          
          <div className="hidden lg:block w-[16px]"></div> {/* Spacer */}
        </header>

        {/* SECTION 1: HERO */}
        <section className="relative w-full h-[100dvh] min-h-[600px] flex flex-col rounded-b-[32px] overflow-hidden pt-32">
          {/* Background Video */}
          <video 
            autoPlay 
            loop 
            muted 
            playsInline 
            className="absolute inset-0 w-full h-full object-cover z-0"
          >
            <source src="https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_045634_e1c98c76-1265-4f5c-882a-4276f2080894.mp4" type="video/mp4" />
          </video>
          
          <div className="relative z-10 flex-1 flex flex-col px-6 sm:px-12 lg:px-24">
            {/* Hero Content */}
            <div className="flex-1 flex flex-col justify-center relative max-w-[780px] lg:ml-32 mt-[-5%]">
              {/* Added Content & Button Container (above the heading) */}
              <div className="mb-8 flex flex-col items-start gap-6 max-w-[600px]">
                <p className="font-mono text-[14px] sm:text-[16px] text-cream/80 leading-relaxed uppercase">
                  Drake is the ultimate environment for multi-agent coding. Run 23+ agents simultaneously, review their pull requests, and orchestrate complex workflows from a single intuitive dashboard.
                </p>
                <Link to="/app" className="liquid-glass px-8 py-4 rounded-full text-neon font-grotesk text-[16px] tracking-widest hover:bg-white/10 transition-colors border border-white/20 uppercase flex items-center gap-2">
                  Launch Workspace
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14"/><path d="m12 5 7 7-7 7"/></svg>
                </Link>
              </div>

              <div className="relative inline-block w-fit">
                <h1 className="font-grotesk text-[40px] sm:text-[60px] md:text-[75px] lg:text-[90px] uppercase leading-[1.05] sm:leading-[1] text-cream">
                  Parallel AI <br />
                  Orchestration Layer
                </h1>
                <div className="font-condiment text-[24px] sm:text-[32px] md:text-[40px] lg:text-[48px] text-neon mix-blend-exclusion opacity-90 absolute right-0 top-[100%] sm:top-auto sm:bottom-0 sm:-right-8 -rotate-1 origin-bottom-right whitespace-nowrap">
                  Agent IDE
                </div>
              </div>
              
              {/* Mobile Socials */}
              <div className="lg:hidden flex items-center justify-center gap-4 mt-12 w-full">
                <a href="#" className="liquid-glass w-[56px] h-[56px] rounded-[1rem] flex items-center justify-center hover:bg-white/10 transition-colors">
                  <Mail className="w-5 h-5 text-cream" />
                </a>
                <a href="#" className="liquid-glass w-[56px] h-[56px] rounded-[1rem] flex items-center justify-center hover:bg-white/10 transition-colors">
                  <TwitterIcon className="w-5 h-5 text-cream" />
                </a>
                <a href="#" className="liquid-glass w-[56px] h-[56px] rounded-[1rem] flex items-center justify-center hover:bg-white/10 transition-colors">
                  <GithubIcon className="w-5 h-5 text-cream" />
                </a>
              </div>
            </div>

            {/* Desktop Socials */}
            <div className="hidden lg:flex flex-col gap-4 absolute right-12 top-[50%] -translate-y-1/2">
              <a href="#" className="liquid-glass w-[56px] h-[56px] rounded-[1rem] flex items-center justify-center hover:bg-white/10 transition-colors">
                <Mail className="w-5 h-5 text-cream" />
              </a>
              <a href="#" className="liquid-glass w-[56px] h-[56px] rounded-[1rem] flex items-center justify-center hover:bg-white/10 transition-colors">
                <TwitterIcon className="w-5 h-5 text-cream" />
              </a>
              <a href="#" className="liquid-glass w-[56px] h-[56px] rounded-[1rem] flex items-center justify-center hover:bg-white/10 transition-colors">
                <GithubIcon className="w-5 h-5 text-cream" />
              </a>
            </div>
          </div>
        </section>

        {/* SECTION 2: ABOUT / INTRO */}
        <section className="relative w-full min-h-[100dvh] flex flex-col py-16 sm:py-24 px-6 sm:px-12 lg:px-24">
          <video 
            autoPlay 
            loop 
            muted 
            playsInline 
            className="absolute inset-0 w-full h-full object-cover z-0"
          >
            <source src="https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_151551_992053d1-3d3e-4b8c-abac-45f22158f411.mp4" type="video/mp4" />
          </video>
          
          <div className="relative z-10 flex-1 flex flex-col justify-between h-full w-full">
            <div className="flex flex-col lg:flex-row lg:justify-between lg:items-start gap-12">
              <div className="relative inline-block w-fit">
                <h2 className="font-grotesk text-[32px] sm:text-[48px] lg:text-[60px] uppercase leading-none text-cream">
                  Hello!<br />
                  I'm Drake
                </h2>
                <div className="font-condiment text-[36px] sm:text-[52px] lg:text-[68px] text-neon mix-blend-exclusion absolute -bottom-6 -right-12 -rotate-2">
                  Drake
                </div>
              </div>
              <p className="font-mono text-[14px] sm:text-[16px] uppercase text-cream max-w-[266px]">
                A shared workspace where agents, terminals, branches, and PRs are supervised from one place.
              </p>
            </div>

          </div>
        </section>

        {/* SECTION 3: AGENTS GRID */}
        <section id="agents" className="w-full bg-[#010828] py-24 px-6 sm:px-12 lg:px-24">
          <div className="flex flex-col lg:flex-row justify-between items-start lg:items-end gap-12 mb-20">
            <h2 className="font-grotesk text-[32px] sm:text-[48px] lg:text-[60px] uppercase leading-none text-cream">
              Network of<br />
              <div className="ml-12 sm:ml-24 lg:ml-32 flex items-baseline gap-4">
                <span className="font-condiment text-neon normal-case">23+</span>
                <span>worker agents</span>
              </div>
            </h2>
            
            <Link to="/app" className="group flex flex-col items-end cursor-pointer">
              <div className="flex items-end gap-2 text-cream">
                <span className="font-grotesk text-[32px] sm:text-[48px] lg:text-[60px] leading-none">START</span>
                <div className="flex flex-col items-start leading-none pb-1">
                  <span className="font-grotesk text-[20px] sm:text-[28px] lg:text-[36px]">WORKSPACE</span>
                  <span className="font-grotesk text-[20px] sm:text-[28px] lg:text-[36px]">NOW</span>
                </div>
              </div>
              <div className="w-full h-[6px] sm:h-[8px] lg:h-[10px] bg-neon mt-2 group-hover:scale-x-110 origin-right transition-transform" />
            </Link>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[
              { status: 'ACTIVE', video: 'https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_053923_22c0a6a5-313c-474c-85ff-3b50d25e944a.mp4', name: 'Claude Code' },
              { status: 'REVIEWING', video: 'https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_054411_511c1b7a-fb2f-42ef-bf6c-32c0b1a06e79.mp4', name: 'Codex' },
              { status: 'STANDBY', video: 'https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_055427_ac7035b5-9f3b-4289-86fc-941b2432317d.mp4', name: 'Aider' }
            ].map((agent, i) => (
              <div key={i} className="liquid-glass rounded-[32px] p-[18px] hover:bg-white/10 transition-colors group flex flex-col">
                <div className="w-full pt-[100%] relative rounded-[24px] overflow-hidden mb-4">
                  <video 
                    autoPlay 
                    loop 
                    muted 
                    playsInline 
                    className="absolute inset-0 w-full h-full object-cover"
                  >
                    <source src={agent.video} type="video/mp4" />
                  </video>
                </div>
                <div className="liquid-glass rounded-[20px] px-5 py-4 flex items-center justify-between">
                  <div className="flex flex-col">
                    <span className="font-mono text-[11px] text-cream/70 uppercase">AGENT: {agent.name}</span>
                    <span className="font-mono text-[16px] text-cream">{agent.status}</span>
                  </div>
                  <Link to="/app" className="w-[48px] h-[48px] rounded-full bg-gradient-to-br from-[#b724ff] to-[#7c3aed] flex items-center justify-center shadow-lg shadow-purple-500/50 hover:scale-110 transition-transform cursor-pointer">
                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-white">
                      <path d="M9 18l6-6-6-6" />
                    </svg>
                  </Link>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* SECTION 4: CTA / FINAL */}
        <section className="relative w-full">
          <video 
            autoPlay 
            loop 
            muted 
            playsInline 
            className="w-full h-auto block"
          >
            <source src="https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_055729_72d66327-b59e-4ae9-bb70-de6ccb5ecdb0.mp4" type="video/mp4" />
          </video>

          <div className="absolute inset-0 flex flex-col justify-center items-end lg:pr-[20%]">
            <div className="relative">
              <div className="font-condiment text-[17px] sm:text-[34px] md:text-[50px] lg:text-[68px] text-neon mix-blend-exclusion absolute -top-8 -left-12 -rotate-3 z-10 whitespace-nowrap">
                Scale your code
              </div>
              <h2 className="font-grotesk text-[16px] sm:text-[32px] md:text-[48px] lg:text-[60px] uppercase text-cream text-right relative z-0">
                <div className="mb-4 sm:mb-8 lg:mb-12">ORCHESTRATE.</div>
                <div>RUN IN PARALLEL.</div>
                <div>RESOLVE CONFLICTS.</div>
                <div>BUILD FASTER.</div>
              </h2>
            </div>
          </div>

          <div className="absolute left-[8%] bottom-[12%] sm:bottom-[16%] lg:bottom-[20%] liquid-glass rounded-[0.5rem] sm:rounded-[1rem] lg:rounded-[1.25rem] flex flex-col">
            <a href="#" className="w-[14vw] sm:w-[14.375rem] md:w-[10.78125rem] lg:w-[16.77rem] py-3 sm:py-4 lg:py-6 flex justify-center items-center hover:bg-white/10 transition-colors border-b border-white/10">
              <Mail className="w-[4vw] h-[4vw] sm:w-6 sm:h-6 lg:w-8 lg:h-8 text-cream" />
            </a>
            <a href="#" className="w-[14vw] sm:w-[14.375rem] md:w-[10.78125rem] lg:w-[16.77rem] py-3 sm:py-4 lg:py-6 flex justify-center items-center hover:bg-white/10 transition-colors border-b border-white/10">
              <TwitterIcon className="w-[4vw] h-[4vw] sm:w-6 sm:h-6 lg:w-8 lg:h-8 text-cream" />
            </a>
            <a href="#" className="w-[14vw] sm:w-[14.375rem] md:w-[10.78125rem] lg:w-[16.77rem] py-3 sm:py-4 lg:py-6 flex justify-center items-center hover:bg-white/10 transition-colors">
              <GithubIcon className="w-[4vw] h-[4vw] sm:w-6 sm:h-6 lg:w-8 lg:h-8 text-cream" />
            </a>
          </div>
        </section>

        {/* SECTION 4.5: FAQ */}
        <section id="faq" className="w-full bg-[#010828] py-24 px-6 sm:px-12 lg:px-24 flex flex-col lg:flex-row items-center gap-16 relative z-10 border-t border-white/10">
          <div className="flex-1 flex flex-col gap-8">
            <h2 className="font-grotesk text-[32px] sm:text-[48px] lg:text-[60px] uppercase leading-none text-cream mb-8">
              Frequently Asked<br />Questions
            </h2>
            
            <div className="flex flex-col gap-6">
              {[
                { q: "How many agents can I run?", a: "You can run up to 23+ worker agents simultaneously in parallel, all coordinated by the Drake orchestration layer." },
                { q: "Does Drake support custom models?", a: "Yes, you can configure custom LLMs for any of your agents to suit your specific workflow." },
                { q: "How are conflicts resolved?", a: "Our platform includes dedicated reviewer agents to handle merge conflicts automatically in the background." },
              ].map((faq, i) => (
                <FaqItem key={i} q={faq.q} a={faq.a} />
              ))}
            </div>
          </div>
          
          <div className="w-full lg:w-[400px] xl:w-[500px] aspect-square liquid-glass rounded-[32px] p-[18px] relative overflow-hidden group">
            <div className="w-full h-full relative rounded-[24px] overflow-hidden">
              <video 
                autoPlay 
                loop 
                muted 
                playsInline 
                className="absolute inset-0 w-full h-full object-cover group-hover:scale-105 transition-transform duration-700"
              >
                <source src="https://d8j0ntlcm91z4.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/hf_20260331_054411_511c1b7a-fb2f-42ef-bf6c-32c0b1a06e79.mp4" type="video/mp4" />
              </video>
              <div className="absolute inset-0 bg-gradient-to-t from-[#010828] to-transparent opacity-60" />
            </div>
          </div>
        </section>

        {/* SECTION 5: FOOTER */}
        <footer className="w-full py-12 px-6 sm:px-12 lg:px-24 border-t border-white/10 flex flex-col sm:flex-row items-center justify-between z-10 relative mt-8">
          <div className="font-grotesk text-[24px] uppercase text-cream tracking-wider mb-4 sm:mb-0">
            Drake
          </div>
          <div className="flex gap-6">
            <a href="#" className="font-mono text-[13px] text-cream/70 hover:text-neon uppercase transition-colors">Terms of Service</a>
            <a href="#" className="font-mono text-[13px] text-cream/70 hover:text-neon uppercase transition-colors">Privacy Policy</a>
          </div>
          <div className="font-mono text-[13px] text-cream/50 mt-4 sm:mt-0 uppercase">
            © 2026 Drake Orchestration Layer. All rights reserved.
          </div>
        </footer>
      </div>
    </div>
  );
}
