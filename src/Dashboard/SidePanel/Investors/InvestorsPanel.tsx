import SidePanel from '../SidePanel';
import InvestorList from './InvestorList';

interface InvestorsPanelProps {
  visible: boolean;
  onClose?: () => void;
}

export default function InvestorsPanel({ visible, onClose }: InvestorsPanelProps) {
  return (
    <SidePanel side="left" visible={visible} onClose={onClose}>
      <InvestorList />
    </SidePanel>
  );
}
