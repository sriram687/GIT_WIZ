import { createFileRoute } from '@tanstack/react-router';
import { OrbisLanding } from '../components/OrbisLanding';

export const Route = createFileRoute('/')({
  component: OrbisLanding,
});
