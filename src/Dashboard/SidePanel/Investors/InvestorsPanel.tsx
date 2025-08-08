import SidePanel from '../SidePanel';
import InvestorList from './InvestorList';
import type { Investor } from './Investor';

interface InvestorsPanelProps {
  visible: boolean;
  onClose?: () => void;
  investors: Investor[];
}

export default function InvestorsPanel({ visible, onClose, investors }: InvestorsPanelProps) {
  return (
    <SidePanel side="left" visible={visible} onClose={onClose}>
      <InvestorList investors={investors} />
    </SidePanel>
  );
}
