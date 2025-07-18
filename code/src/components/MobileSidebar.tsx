'use client';

import { useState, type ReactElement, type ReactNode } from 'react';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/24/outline';
import { usePathname } from 'next/navigation';

interface MenuItem {
    name: string;
    icon: ReactElement;
    href: string;
}

interface SidebarProps {
  menuItems: MenuItem[];
  accountDropdown?: ReactNode;
}

export default function MobileSidebar({ menuItems, accountDropdown } : SidebarProps) {
  const [isMobileSidebarOpen, setIsMobileSidebarOpen] = useState(false);

  const pathname = usePathname();
  
  // Find the menu item that matches the current path
  const activeItem = menuItems.find(item => pathname.startsWith(item.href))?.name || menuItems[0]?.name;

  return (
    <>
      {/* Mobile Hamburger Button */}
      <button
        onClick={() => setIsMobileSidebarOpen(true)}
        className="lg:hidden fixed top-4 right-4 z-50 p-2 text-white 
                  transition-all duration-300"
        aria-label="Open sidebar"
      >
        <Bars3Icon className="h-6 w-6 scale-150 cursor-pointer" />
      </button>

      {/* Mobile Sidebar - Fullscreen */}
      {isMobileSidebarOpen && (
        <div className="lg:hidden fixed inset-0 z-50">
          <aside
            className={`absolute inset-0 bg-neutral-800 transform ${
              isMobileSidebarOpen ? 'translate-x-0' : '-translate-x-full'
            } transition-transform duration-300 ease-in-out`}
          >
            <div className="relative flex flex-col h-full max-h-full">
              <header className="p-4 mt-1 ml-1 flex justify-between items-center gap-x-2">
                <a className="flex-none font-semibold text-xl" href="#" aria-label="Brand">
                  Ekstra
                </a>
                <button
                  onClick={() => setIsMobileSidebarOpen(false)}
                  className="p-1 rounded-md hover:bg-neutral-700 hover:cursor-pointer"
                  aria-label="Close sidebar"
                >
                  <XMarkIcon className="h-6 w-6 cursor-pointer" />
                </button>
              </header>
              
              <nav className="flex-1 overflow-y-auto pt-1.5">
                <div className="p-2 w-full flex flex-col flex-wrap">
                  <ul className="space-y-1">
                    {menuItems.map((item) => (
                      <li key={item.name}>
                        <a
                          className={`w-full flex items-center gap-x-3.5 p-2.5 text-base rounded-lg hover:bg-neutral-700 ${
                            activeItem === item.name ? "bg-neutral-700" : ""
                          }`}
                          href={item.href}
                        >
                          {item.icon}
                          {item.name}
                        </a>
                      </li>
                    ))}
                  </ul>
                </div>
              </nav>
              
              {/* Mobile Account Info Section */}
              <footer className="p-4 border-t border-neutral-700">
                {accountDropdown}
              </footer>
            </div>
          </aside>
        </div>
      )}
    </>
  );
}