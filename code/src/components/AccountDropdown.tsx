"use client"

import { useRef } from 'react';
import { signOut } from "next-auth/react";
import { useSession } from "next-auth/react";
import { useClickOutside } from '~/hooks/useClickOutside';
import { useState } from 'react';

export default function AccountDropdown() {
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const buttonRef = useRef<HTMLButtonElement>(null);

  const { data: session } = useSession();
  const userName = session?.user?.name;
  const userImage = session?.user?.image ?? "https://www.gravatar.com/avatar/?d=mp";
  const creatorPageId = session?.user?.creatorPageId;

  const dropDownItemStyle = 
    "block px-4 py-2 text-sm text-neutral-200 hover:bg-gray-100 hover:bg-neutral-700 w-full text-left cursor-pointer"

  useClickOutside({
    refs: [dropdownRef, buttonRef],
    handler: () => setIsDropdownOpen(false),
    isActive: isDropdownOpen,
    eventType: 'mousedown'
  });

  return (
    <div className="relative">
      <button 
        ref={buttonRef}
        onClick={(e) => {
          e.stopPropagation();
          setIsDropdownOpen(!isDropdownOpen);
        }}
        className="w-full flex items-center gap-x-3 p-2 rounded-lg hover:bg-neutral-700 transition-colors duration-200 cursor-pointer"
        aria-expanded={isDropdownOpen}
        aria-label="Account menu"
      >
        <img 
          className="w-8 h-8 rounded-full" 
          src={userImage} 
          alt="User avatar"
          referrerPolicy="no-referrer"
        />
        <div className="text-left">
          <p className="text-sm font-medium text-neutral-200">{userName}</p>
        </div>
        <svg 
          className={`ml-auto h-4 w-4 text-neutral-400 transition-transform duration-200 ${
            isDropdownOpen ? 'rotate-180' : ''
          }`}
          xmlns="http://www.w3.org/2000/svg" 
          viewBox="0 0 24 24" 
          fill="none" 
          stroke="currentColor" 
          strokeWidth="2" 
          strokeLinecap="round" 
          strokeLinejoin="round"
        >
          <path d="m7 15 5 5 5-5"/>
          <path d="m7 9 5-5 5 5"/>
        </svg>
      </button>
      
      {isDropdownOpen && (
        <div ref={dropdownRef} className="absolute bottom-full left-0 mb-2 w-full bg-neutral-800 rounded-lg shadow-lg border border-neutral-700 overflow-hidden z-10">
          {/* Options depending on user creator status */}
            {!creatorPageId && (
              <a 
                href="/creatorsetup" className={dropDownItemStyle}>
                Create on Ekstra
              </a>
            )}
            {creatorPageId && (
              <a 
                href="/creator/dashboard/home"
                className={dropDownItemStyle}
              >
                Creator Dashboard
              </a>
            )}

          <a 
            href="#" 
            className={dropDownItemStyle}
          >
            Profile
          </a>
          <button
            onClick={() => signOut({ redirectTo: "/" })}
            className={dropDownItemStyle}
          >
            Sign out
          </button>
        </div>
      )}
    </div>
  );
}