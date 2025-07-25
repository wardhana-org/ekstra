interface Membership {
  id: string,
  name: string,
  price: number,
  description: string,
}

export const plans:Membership[] = [
  { 
    id: '1', 
    name: 'Bronze Tier', 
    price: 5,
    description: 'Basic access to exclusive content and community features' 
  },
  { 
    id: '2', 
    name: 'Silver Tier', 
    price: 10,
    description: 'Early access to new content plus all Bronze benefits' 
  },
  { 
    id: '3', 
    name: 'Gold Tier', 
    price: 25,
    description: 'VIP access with exclusive perks and all lower tier benefits' 
  }
];