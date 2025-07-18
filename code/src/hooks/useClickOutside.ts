"use client"

import { useEffect, type RefObject } from 'react';

type ClickOutsideConfig = {
  refs: Array<RefObject<HTMLElement | null>>; // Allow null here
  handler: () => void;
  isActive?: boolean;
  eventType?: keyof DocumentEventMap;
};

export function useClickOutside({
  refs,
  handler,
  isActive = true,
  eventType = 'mousedown'
}: ClickOutsideConfig) {
  useEffect(() => {
    if (!isActive) return;

    const handleEvent = (event: Event) => {
      const isOutside = refs.every(ref => {
        return ref.current && !ref.current.contains(event.target as Node);
      });

      if (isOutside) {
        handler();
      }
    };

    document.addEventListener(eventType, handleEvent);
    return () => {
      document.removeEventListener(eventType, handleEvent);
    };
  }, [refs, handler, isActive, eventType]);
}