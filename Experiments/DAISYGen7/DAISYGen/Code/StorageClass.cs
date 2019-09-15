using System;
using System.Collections.Generic;
using Gtk;

namespace DAISYGen
{
	// class for containing of the responsive components:
	// associated widget, is it switchable, current state of the widget (on/off)
	
	public class Storage
	{
		public Widget widget;
		public bool switchable;
		public bool current_state; 
		public Storage(Widget w, bool sw=false, bool cs=false)
		{
			widget = w;
			switchable = sw;
			current_state = cs;
		}
	}
}
