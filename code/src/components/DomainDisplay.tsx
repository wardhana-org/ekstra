import { useState, useEffect } from "react";

export default function DomainDisplay() {
  const [domain, setDomain] = useState('');

  useEffect(() => {
    setDomain(window.location.hostname);
  }, []);

  return <span>{domain}</span>;
}