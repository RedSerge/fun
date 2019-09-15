// Main GUI and Event Handling Code file for DAISYGen Project 

using System;
using System.IO;
using System.Collections.Generic;
using System.Runtime.Serialization.Formatters.Binary;
using Gtk;
using DAISYGen;

public partial class MainWindow: Gtk.Window
{	
	//List of Audio Tracks (L.A.T.)
	public static List<AudioTracks> lat = new List<AudioTracks> ();
	//Index of Current Item selected on the panel (Audio, Timestamp)
	public static int Audio_GUI_Index = -1;
	public static int Stamp_GUI_Index = -1;
	//Storage of Custom Image Controls
	List<Storage> ControlStorage = new List<Storage> ();
	//Flag representing the current mode of interaction with timestamp widgets (search or modification)
	bool ActivationAsSearch=false;
	//The lists of Dynamic Widgets
	List<Widget> Dynamic_Widgets_Audio;
	List<Widget> Dynamic_Widgets_Timestamps;
	List<Widget> Dynamic_Widgets_TimestampsOffset;
	List<Widget> Dynamic_Widgets_TimestampsCaption;
	//Flag controlling the behaviour of Media Player component in regards of L.A.T.
	bool reset_duration = false;
	//Current sound volume level (while sound is muted)
	int storedvol = 100;
	//Measure of occupancy of the custom components featuring the duration and volume scale  
	double soundratio = 1;
	double playerratio = 0;
	
	// Function expands given number of seconds to the timespan string representation
	string SecondsToFullTime(double sec)
	{
		var ts = new TimeSpan (0, 0, (int)sec);
		return String.Join (":", ts.Hours.ToString ("00"), ts.Minutes.ToString ("00"), ts.Seconds.ToString ("00"));
	}
	
	//Function dynamically updates the player state
	public bool PlayerTimer()
	{
		//Check if any specific audio track is selected
		if (Audio_GUI_Index > 0) {
			//1st label describing position
			double pos = MainClass.wmp.Ctlcontrols.currentPosition;
			label1.Text = SecondsToFullTime (pos);
			//2nd label describing duration
			double dur = MainClass.wmp.currentMedia.duration;
			label2.Text = SecondsToFullTime (dur);
			//Redraw the whole bar between labels
			playerratio = pos / dur;
			wavebar.QueueDraw ();
		} else {
			// No audio track is selected? Reset labels, then
			if (label2.Text != "00:00:00") {
				label1.Text = label2.Text = "00:00:00";
			}
		}
		//Keep going
		return true;
	}

	//Function for changing every GTK widget color state
	public void SetWidgetColor(Widget widget, string color_string)
	{
		Gdk.Color c = new Gdk.Color ();
		Gdk.Color.Parse (color_string, ref c);
		widget.ModifyBase (StateType.Normal, c);
		widget.ModifyBase (StateType.Active, c);
		widget.ModifyBase (StateType.Insensitive, c);
		widget.ModifyBase (StateType.Prelight, c);
		widget.ModifyBase (StateType.Selected, c);
		widget.ModifyBg (StateType.Normal, c);
		widget.ModifyBg (StateType.Active, c);
		widget.ModifyBg (StateType.Insensitive, c);
		widget.ModifyBg (StateType.Prelight, c);
		widget.ModifyBg (StateType.Selected, c);
		widget.ModifyFg (StateType.Normal, c);
		widget.ModifyFg (StateType.Active, c);
		widget.ModifyFg (StateType.Insensitive, c);
		widget.ModifyFg (StateType.Prelight, c);
		widget.ModifyFg (StateType.Selected, c);		
	}

	// Function converts R, G, B channel values into #RRGGBB hexademical string 
	public string RGB2HEX(int r, int g, int b)
	{
		string rc = Convert.ToString (r, 16).PadLeft (2, '0').ToUpper ();
		string gc = Convert.ToString (g, 16).PadLeft (2, '0').ToUpper ();
		string bc = Convert.ToString (b, 16).PadLeft (2, '0').ToUpper ();
		return "#" + rc + gc + bc;
	}

	//Function converts given color channel value to the real range [0.0; 1.0]
	public double ColorChan(int chan)
	{
		return chan / 255.0;
	}
	
	//Function fills the control storage based on its name
	void ImagesAsControls()
	{
		ControlStorage.Clear ();
		ControlStorage.Add (new Storage (skip_previous));
		//true,true <==> Switchable, on
		ControlStorage.Add (new Storage (play, true, true));
		ControlStorage.Add (new Storage (skip1));
		ControlStorage.Add (new Storage (stop));
		//Switchable, on
		ControlStorage.Add (new Storage (sound, true, true));
		ControlStorage.Add (new Storage (info));
		ControlStorage.Add (new Storage (create));
		ControlStorage.Add (new Storage (open));
		ControlStorage.Add (new Storage (save));
		ControlStorage.Add (new Storage (delete_file));
		//"new" is a reserved word; and yet it can be used as a name this way:
		ControlStorage.Add (new Storage (@new));
		ControlStorage.Add (new Storage (up));
		ControlStorage.Add (new Storage (down));
		ControlStorage.Add (new Storage (delete));
		ControlStorage.Add (new Storage (write));
		//Switchable, off
		ControlStorage.Add (new Storage (activation, true));
		ControlStorage.Add (new Storage (daisy));
		ControlStorage.Add (new Storage (window_minimize));
		ControlStorage.Add (new Storage (window_close));
	}
	
	//Function resets all vital properties and sets necessary handlers for the elements of the Control Storage
	void ActivateImages()
	{
		foreach (Storage s in ControlStorage) {
			//Visual representation
			var imagecontrol = (Image)s.widget;
			EventBox eb = (EventBox)imagecontrol.Parent;
			Alignment align = (Alignment)eb.Parent;
			align.Xscale = align.Yscale = imagecontrol.Xalign = imagecontrol.Yalign = 0;
			eb.VisibleWindow = false;
			//Handlers configuration 
			eb.Parent.Events = play.Events;
			eb.ButtonPressEvent += new ButtonPressEventHandler (HandleImages);
			eb.ButtonReleaseEvent += new ButtonReleaseEventHandler (HandleReleases);
			eb.EnterNotifyEvent += new EnterNotifyEventHandler (HandleEnters);
			eb.LeaveNotifyEvent += new LeaveNotifyEventHandler (HandleAutumn);
			//Redraw Image
			DrawActiveImage (eb, 0);
		}
	}

	//Add audiotrack dialog and corresponding procedure 
	public bool AddAudioDialog()
	{
		//Checking the capacity overflow
		if (lat.Count >= Dynamic_Widgets_Audio.Count - 2) {				
			Alert.Show (this, "The audiotrack limit has been reached");
			return false;
		}
		//Prepare dialog and extension filters
		var temporarydialog = new Gtk.FileChooserDialog ("", this, FileChooserAction.Open, "Cancel", ResponseType.Cancel, "Open", ResponseType.Accept);
		temporarydialog.SetCurrentFolder (@".");
		var filefilter_all = new FileFilter ();
		filefilter_all.Name = "All Files (*.*)";
		filefilter_all.AddPattern ("*.*");
		var filefilter_wave = new FileFilter ();
		filefilter_wave.Name = "Standard Audio Files (*.mp3;*.wav)";
		filefilter_wave.AddPattern ("*.mp3");
		filefilter_wave.AddPattern ("*.wav");
		temporarydialog.AddFilter (filefilter_wave);
		temporarydialog.AddFilter (filefilter_all);
		//Check if file has been selected
		bool result = temporarydialog.Run () == (int)ResponseType.Accept;
		if (result) {
			//Add selected file and refresh corresponding panel
			lat.Add (new AudioTracks (temporarydialog.Filename));
			Refresh_Audio_Widgets (lat.Count);
		}
		temporarydialog.Destroy ();
		return result;
	}

	//Function shows the custom dialog based on provided action description
	string ShowCustomDialog(FileChooserAction fca, string action)
	{
		var temporarydialog = new Gtk.FileChooserDialog ("", this, fca, "Cancel", ResponseType.Cancel, action, ResponseType.Accept);
		temporarydialog.SetCurrentFolder (@".");
		//Only *.pro files are supposed to be supported
		var filefilter = new FileFilter ();
		filefilter.Name = "Project Files (*.pro)";
		filefilter.AddPattern ("*.pro");
		temporarydialog.AddFilter (filefilter);
		//Returns empty string on cancellation, filename otherwise. 
		string export = string.Empty;
		if (temporarydialog.Run () == (int)ResponseType.Accept) {			
			export = temporarydialog.Filename;
			//Check if extension is presented, enforce its presence
			string export_ext = System.IO.Path.GetExtension (export); 
			export += export_ext == ".pro" ? string.Empty : ".pro";
		}
		temporarydialog.Destroy ();
		return export;
	}

	//Save dialog and corresponding procedure
	void Show_SaveDialog()
	{
		//Binary serialization of list of audio tracks with export to the file with given filename
		var filename = ShowCustomDialog (FileChooserAction.Save, "Save");
		if (filename == string.Empty)
			return;
		var bf = new BinaryFormatter ();
		var ms = new MemoryStream ();
		bf.Serialize (ms, lat);
		File.WriteAllBytes (filename, ms.ToArray ());
		//Update the Project Name in the title
		MainTitle.Text = System.IO.Path.GetFileNameWithoutExtension (filename);
	}

	//Open dialog and corresponding procedure
	void Show_OpenDialog()
	{	
		//Export the list of audio tracks from corresponding binary serialized file with given filename
		var filename = ShowCustomDialog (FileChooserAction.Open, "Open");
		if (filename == string.Empty)
			return;	
		var bf = new BinaryFormatter ();
		var ms = new MemoryStream ();
		var content = File.ReadAllBytes (filename);
		ms.Write (content, 0, content.Length);
		ms.Position = 0;
		lat = bf.Deserialize (ms) as List<AudioTracks>;
		//Update and restore the GUI and subsystems
		Stamp_GUI_Index = Audio_GUI_Index = -1;
		Refresh_Audio_Widgets ();
		Refresh_Timestamp_Widgets ();
		MainClass.wmp.Ctlcontrols.stop();
		MainClass.wmp.currentPlaylist.clear();
		MainTitle.Text = System.IO.Path.GetFileNameWithoutExtension (filename);
	}

	//Function launches corresponding actions based on activated widget
	void UnifiedActionControl(Storage storedwidget, string widgetname)
	{
		switch (widgetname) {
		//Convert to DAISY format
		case "daisy":
			ExportClass.Extract (this, lat);
			break;
		//Delete current audiotrack
		case "delete_file":
			if (Audio_GUI_Index > 0) {
				int index = Audio_GUI_Index - 1;							
				lat.RemoveAt (index);
				if (index >= lat.Count)
					Audio_GUI_Index--;
				Refresh_Audio_Widgets ();
				Stamp_GUI_Index = -1;
				Refresh_Timestamp_Widgets ();
			}
			break;
		//Delete the whole project from memory
		case "create":
			StopPlay ();
			Audio_GUI_Index = -1;
			Stamp_GUI_Index = -1;
			lat.Clear ();
			Refresh_Audio_Widgets ();
			Refresh_Timestamp_Widgets ();
			MainTitle.Text = "[*No project*]";
			break;
		//Move track to the previous position, if available
		case "skip_previous":
			if (Audio_GUI_Index > 0) {
				int index = Audio_GUI_Index - 1;
				if (index > 0) {
					var lost_element = lat [index];
					lat.RemoveAt (index);
					lat.Insert (index - 1, lost_element);
					Audio_GUI_Index--;
					Refresh_Audio_Widgets (Audio_GUI_Index, true);
				}
			}
			break;
		//Move track to the next position, if available
		case "skip1":
			if (Audio_GUI_Index > 0) {
				int index = Audio_GUI_Index - 1;
				if (index < lat.Count - 1) {
					var lost_element = lat [index];
					lat.RemoveAt (index);
					lat.Insert (index + 1, lost_element);
					Audio_GUI_Index++;
					Refresh_Audio_Widgets (Audio_GUI_Index, true);
				}
			}
			break;
		//Show Open Dialog
		case "open":
			Show_OpenDialog ();
			break;
		//Show Save Dialog
		case "save":
			Show_SaveDialog ();
			break;
		//Minimize window
		case "window_minimize":
			this.Iconify ();
			break;
		//Close window
		case "window_close":
			OnDeleteEvent (this, null);
			break;
		//Mute sound
		case "sound":
			if (!storedwidget.current_state) {
				storedvol = MainClass.wmp.settings.volume;
				MainClass.wmp.settings.volume = 0;
			} else {
				MainClass.wmp.settings.volume = storedvol;
				storedvol = -1;
			}
			break;
		//About window
		case "info":
			Alert.Show (this, "About", true);
			break;
		//Play or pause
		case "play":
			if (Audio_GUI_Index > 0)
			if (!storedwidget.current_state) {
				MainClass.wmp.Ctlcontrols.play ();
			} else {
				MainClass.wmp.Ctlcontrols.pause ();
			}
			break;
		//Stop playing
		case "stop":
			MainClass.wmp.Ctlcontrols.stop ();
			//Find "play" button and force default state
			var otherstoredwidget = ControlStorage.Find (y => y.widget.Name == "play");
			otherstoredwidget.current_state = true;
			((Image)otherstoredwidget.widget).Pixbuf = Gdk.Pixbuf.LoadFromResource ("DAISYGen.Buttons.play_on_usual.png");
			break;
		//Add new timestamp for the selected audiotrack
		case "new":
			if (Audio_GUI_Index > 0) {
				AddDaisyTime (lat [Audio_GUI_Index - 1].Timestamps);
				Refresh_Timestamp_Widgets ();
			}
			break;
		//Move selected timestamp higher up the hierarchy 
		case "up":
			if (Audio_GUI_Index > 0 && Stamp_GUI_Index > 0) {				
				ChangeHierarchyLevel (lat [Audio_GUI_Index - 1].Timestamps, Stamp_GUI_Index, -1);
				Refresh_Timestamp_Widgets ();
			}
			break;
		//Move selected timestamp lower down the hierarchy
		case "down":
			if (Audio_GUI_Index > 0 && Stamp_GUI_Index > 0) {
				if (lat [Audio_GUI_Index - 1].Timestamps [Stamp_GUI_Index].position < lat [Audio_GUI_Index - 1].Timestamps [Stamp_GUI_Index - 1].position + 1)
					ChangeHierarchyLevel (lat [Audio_GUI_Index - 1].Timestamps, Stamp_GUI_Index, 1);			
				Refresh_Timestamp_Widgets ();
			}
			break;
		//Erase timestamp
		case "delete":
			if (Audio_GUI_Index > 0 && Stamp_GUI_Index > 0) {
				EraseStamp (lat [Audio_GUI_Index - 1].Timestamps, Stamp_GUI_Index);
				Refresh_Timestamp_Widgets ();
			}
			break;
		//Switch activation mode
		case "activation":
			ActivationAsSearch = storedwidget.current_state;
			break;
		//Change current timestamp title
		case "write":
			if (Audio_GUI_Index > 0 && Stamp_GUI_Index > -1) {		
				lat [Audio_GUI_Index - 1].Timestamps [Stamp_GUI_Index].title = textbox.Text;
				Refresh_Timestamp_Widgets ();
			}
			break;
		}
	}
	
	//Function adds timestamp for current audiotrack
	void AddDaisyTime(List<DaisyTime>Stamps, bool zero = false)
	{
		//Discover current position
		double position = zero ? 0 : MainClass.wmp.Ctlcontrols.currentPosition;
		string position_to_seconds = position.ToString ("0.000");
		//Can't add timestamp if it already exists
		if (Stamps.FindIndex (y => y.sec.ToString ("0.000") == position_to_seconds) != -1)
			return;
		//Store the timestamp and sort according to the timeline
		Stamps.Add (new DaisyTime (position, position_to_seconds + "s", -1));
		Stamps.Sort (MainClass.time_comparer);
		//Adjust the hierarchy level for every provided timestamp
		int j = Stamps.Count;
		for (int i = 0; i < j; ++i) {
			if (Stamps [i].position == -1) {
				Stamps [i].position = (i == 0) ? 0 : Stamps [i - 1].position;
			}
		}
		Stamp_GUI_Index = Stamps.FindIndex (y => y.sec.ToString ("0.000") == position_to_seconds);
	}

	//Function draws image properly for given active widget in certain state
	void DrawActiveImage(object o, int state)
	{
		//Extract Image Widget from corresponding Event Box
		Image widget = (Image)((EventBox)o).Child;
		//Find the Widget in the Control Storage
		var storedwidget = ControlStorage.Find (y => y.widget == widget);
		if (storedwidget == null)
			return;
		//Trying to guess a part of necessary resource name based on the value of the state
		string ResourceName;
		switch (state) {
		case 1:
			ResourceName = "_hover.png";
			break;
		case 2:
			ResourceName = "_click.png";		
			break;
		case 3:
			ResourceName = "_hover.png";
			//Additional correction of switchable objects
			if (storedwidget.switchable)
				storedwidget.current_state = !storedwidget.current_state;			
			break;
		default:
			ResourceName = "_usual.png";
			break;
		}
		ResourceName = "DAISYGen.Buttons." + widget.Name + (!storedwidget.switchable ? "" : storedwidget.current_state ? "_on" : "_off") + ResourceName;
		//Loading the resource
		widget.Pixbuf = Gdk.Pixbuf.LoadFromResource (ResourceName);
		//Additional manipulation with switchable objects in case of corresponding event
		if (state == 3)
			UnifiedActionControl (storedwidget, widget.Name);
	}

	//Group of handlers for custom image control events,
	//setting reactions to "Press", "Release", "Enter" and "Leave" controls.

	protected void HandleImages (object o, ButtonPressEventArgs args)
	{
		DrawActiveImage (o, 2);
	}

	protected void HandleReleases (object o, ButtonReleaseEventArgs args)
	{
		DrawActiveImage (o, 3);
	}

	protected void HandleEnters (object o, EnterNotifyEventArgs args)
	{
		DrawActiveImage (o, 1);
	}

	//You know... Autumn... Leaves... Handle leaves... *choked laugh* Nevermind.
	protected void HandleAutumn (object o, LeaveNotifyEventArgs args)
	{
		DrawActiveImage (o, 0);
	}

	//Function starts playing the current audiotrack
	public void BeginPlay(string path_to_audiotrack)
	{
		MainClass.wmp.URL = path_to_audiotrack;
		var widget = ControlStorage.Find (y => y.widget.Name == "play");
		widget.current_state = false;
		((Image)widget.widget).Pixbuf = Gdk.Pixbuf.LoadFromResource ("DAISYGen.Buttons.play_off_usual.png");
	}
	
	//Function stops playing the audiotrack and resets the corresponding custom control
	public void StopPlay()
	{
		MainClass.wmp.Ctlcontrols.stop ();
		var widget = ControlStorage.Find (y => y.widget.Name == "play");
		widget.current_state = false;
		((Image)widget.widget).Pixbuf = Gdk.Pixbuf.LoadFromResource ("DAISYGen.Buttons.play_on_usual.png");
		playerratio = 0;
		wavebar.QueueDraw ();
	}

	//Function refreshes the list of audiotracks
	public void Refresh_Audio_Widgets(int index = -1, bool ignore = false)
	{
		//Update index of selected audiotrack
		if (!ignore) {
			if (index > 0)
				Audio_GUI_Index = index;
		}
		//Hide widgets...
		foreach (Widget widget in Dynamic_Widgets_Audio)
			widget.Visible = false;
		//...Except for the special 1st
		Dynamic_Widgets_Audio[0].Visible = true;
		//Total count of active widgets
		var limit = Math.Min (lat.Count, Dynamic_Widgets_Audio.Count - 2);		
		for (int i = 1; i <= limit; ++i) {
			//Set active or inactive(base) color based on index
			SetWidgetColor (Dynamic_Widgets_Audio [i], RGB2HEX (161, 161, i == Audio_GUI_Index ? 255 : 161));
			//Set appropriate name for the widget
			var inner_text = (Label)(((EventBox)Dynamic_Widgets_Audio[i]).Child);
			string fullname = lat[i - 1].path;
			string shortname = System.IO.Path.GetFileNameWithoutExtension (fullname);
			//Bold font style for selected one
			inner_text.Text = (i == Audio_GUI_Index) ? "<b>" + shortname + "</b>" : shortname;
			inner_text.TooltipText = fullname;
			//HTML on
			inner_text.UseMarkup = true;
			//Show the widget
			Dynamic_Widgets_Audio[i].Visible = true;
		}
		//Last element without corresponding audiotrack is used to add new tracks 
		Dynamic_Widgets_Audio[limit + 1].Visible = true;
		var last_widget = (Label)(((EventBox)Dynamic_Widgets_Audio[limit + 1]).Child);
		SetWidgetColor (Dynamic_Widgets_Audio[limit + 1], RGB2HEX (161, 161, 161));
		last_widget.TooltipText = "Add an element";
		last_widget.Text = "* New audiotrack *";
		//Set the appropriate state of the Windows Media Player
		if (!ignore) {
			if (Audio_GUI_Index > 0)
				BeginPlay (lat[Audio_GUI_Index - 1].path);
			else
				StopPlay ();
		}
	}

	//Function refreshes the list of corresponding timestamps
	public void Refresh_Timestamp_Widgets(int index=-1)
	{
		//Update index		
		if (index > 0)
			Stamp_GUI_Index = index;
		//Hide all widgets but the 1st one
		foreach (Widget widget in Dynamic_Widgets_Timestamps)
			widget.Visible = false;
		Dynamic_Widgets_Timestamps [0].Visible = true;
		//If no audiotrack is selected, stop
		if (Audio_GUI_Index < 1)
			return;
		//Show widgets, set hierarchy offsets and gradient colors, set highlight color
		var stamps = lat [Audio_GUI_Index - 1].Timestamps;
		int j = stamps.Count;
		for (int i = 0; i < j; ++i) {
			var current_widget = Dynamic_Widgets_Timestamps[i + 1];
			((Label)Dynamic_Widgets_TimestampsCaption[i]).Text = "[" + SecondsToFullTime (stamps[i].sec) + "] '" + stamps[i].title + "'";
			int current_widget_index = stamps[i].position;
			int base_stamp_color = 115 + 15 * current_widget_index;
			int higlight_stamp_color = i == Stamp_GUI_Index ? 255 : base_stamp_color;
			SetWidgetColor (current_widget, RGB2HEX (base_stamp_color, base_stamp_color, higlight_stamp_color));
			((EventBox)Dynamic_Widgets_TimestampsOffset[i]).WidthRequest = 100 * current_widget_index;
			current_widget.Visible = true;
		}
	}

	//Erase the widget and its inner timestamps
	void EraseStamp(List<DaisyTime>Stamps, int pos)
	{
		//Check if given position is valid
		if ((pos < 0) || (pos >= Stamps.Count))
			return;
		//Store hierarchy level
		int hierarchy_level = Stamps[pos].position;
		//Erase every inner widget (delete a widget while its hierarchy level is less or equal to the given level) 
		bool flag = true;
		do
		{
			Stamps.RemoveAt(pos);
			if ((pos == Stamps.Count) || (Stamps[pos].position <= hierarchy_level))
				flag = false;
		} while (flag);
	}

	//Change the hierarchy level for timestamps
	void ChangeHierarchyLevel(List<DaisyTime>Stamps, int pos, int dir = -1)
	{
		//Check if given position is valid
		int count = Stamps.Count;
		if ((pos < 0) || (pos >= count))
			return;
		//Detect current level
		int current_level = Stamps [pos].position;
		//Detect level change in given direction
		int changed_level = current_level + dir;
		//Hierarchy level belongs to integer range [0;5]
		changed_level = changed_level < 0 ? 0 : changed_level > 5 ? 5 : changed_level;		
		//Change each element hierarchy under provided one
		for (int i = pos + 1; i < count; ++i) {
			//stop changes outside of hierarchy range of element at given position
			if (Stamps [i].position <= current_level)
				break;
			Stamps [i].position += dir;
			if (Stamps [i].position < 0)
				Stamps [i].position = 0;
			if (Stamps [i].position > 5)
				Stamps [i].position = 5;
		}
		//Change provided element
		Stamps [pos].position = changed_level;
	}

	//2 next functions fill left and right side of the workarea, correspondingly 
	public List<Widget> FillLeftScrolledWindow(ScrolledWindow scroll0)
	{
		//Prepare GUI and logical structures for storing the elements
		var result = new List<Widget> ();
		var event_wrapper = new EventBox ();
		var vertical_container = new VBox ();
		vertical_container.Spacing = 10;
		event_wrapper.Add (vertical_container);
		SetWidgetColor (event_wrapper, RGB2HEX (201, 201, 201));
		//Make a lot of slots for the elements
		for(int i = 0; i < 502; i++)
		{
			var inner_event = new EventBox ();
			if (i == 0)
				//Special first element
				SetWidgetColor (inner_event, RGB2HEX (201, 201, 201));
			else {
				//Ordinary element with HTML-supporting label
				SetWidgetColor (inner_event, RGB2HEX (161, 161, 161));
				var inner_label = new Label ("");
				inner_label.UseMarkup = true;
				inner_event.Add (inner_label);
			}
			//Set the default properties for the slot
			inner_event.HeightRequest = i == 0 ? 1 : 70;
			inner_event.Name = "F1#" + result.Count.ToString ();
			inner_event.ButtonReleaseEvent += new ButtonReleaseEventHandler (AudioTrackClick);
			inner_event.HasTooltip = true;
			//Add the slot to the container and the resulting structure
			vertical_container.PackStart (inner_event, false, false, 0);
			result.Add (inner_event);
		}
		//Connect the system of elements to the scrollable container
		scroll0.AddWithViewport (event_wrapper);
		scroll0.ShowAll ();
		scroll0.WindowPlacement = CornerType.TopRight;
		scroll0.WindowPlacementSet = true;
		//Hide the elements, as they are not initialized properly yet
		foreach (Widget widget in result)
			widget.Visible = false;
		//Return the list of added elements
		return result;
	}

	public List<Widget> FillRightScrolledWindow(ScrolledWindow scroll0)
	{
		var result = new List<Widget> ();
		//Must have two groups of widgets for this panel: captions and offsets
		Dynamic_Widgets_TimestampsCaption = new List<Widget> ();
		Dynamic_Widgets_TimestampsOffset = new List<Widget> ();
		var event_wrapper = new EventBox ();
		SetWidgetColor (event_wrapper, RGB2HEX (201, 201, 201));
		var vertical_container = new VBox ();
		vertical_container.Spacing = 10;
		event_wrapper.Add (vertical_container);
		for(int i = 0; i < 501; i++)
		{
			var inner_event = new EventBox ();
			if (i == 0)
				SetWidgetColor (inner_event, RGB2HEX (201, 201, 201));
			else {				
				SetWidgetColor (inner_event, RGB2HEX (115, 115, 115));
				var offset_keeper = new HBox ();
				var offset_setter = new EventBox ();
				Dynamic_Widgets_TimestampsOffset.Add (offset_setter);
				SetWidgetColor (offset_setter, RGB2HEX (201, 201, 201));
				//Default offset is zero (may be changed during the runtime based on the hierarchy level)
				offset_setter.WidthRequest = 0;
				var inner_label = new Label ("");
				Dynamic_Widgets_TimestampsCaption.Add (inner_label);
				inner_label.UseMarkup = true;
				//Allow the text to be wrapped
				inner_label.Wrap = true;
				offset_keeper.PackStart (offset_setter, false, false, 0);
				offset_keeper.PackStart (inner_label, true, true, 0);
				inner_event.Add (offset_keeper);
			}
			inner_event.HeightRequest = i == 0 ? 1 : 70;
			inner_event.Name = "F2#" + result.Count.ToString ();
			inner_event.ButtonReleaseEvent += new ButtonReleaseEventHandler (TimeStampClick);
			inner_event.HasTooltip = true;
			vertical_container.PackStart (inner_event, false, false, 0);
			result.Add (inner_event);
		}
		scroll0.AddWithViewport (event_wrapper);
		scroll0.ShowAll ();
		foreach (Widget widget in result)
			widget.Visible = false;
		return result;
	}

	//This handler controls the change of the Windows Media Player state
	public void AxWindowsMediaPlayer1StatusChange(object sender, EventArgs e)
	{
		//Get current track duration
		var duration = MainClass.wmp.currentMedia.duration;
		//Check if any specific audio track is selected
		if (Audio_GUI_Index > 0) {
			//Check if the duration of the track has been initialized
			if (duration != 0) {
				if (reset_duration) {
					lat [Audio_GUI_Index - 1].length = duration;
					reset_duration = false;
				}
			} else {
				reset_duration = true;
			}
		}
	}

//REMAINS UNCHANGED YET **********************

	//Handler for clicking the audio tracks
	public void AudioTrackClick(object o,ButtonReleaseEventArgs args)
	{
		//Get index of selected item
		Audio_GUI_Index = int.Parse (((Widget)o).Name.Split ('#') [1]);
		//Try to add new audio track if the corresponding special item is activated
		if (Audio_GUI_Index > lat.Count) {
			Audio_GUI_Index = -1;
			Stamp_GUI_Index = -1;
			if (!AddAudioDialog ()) {
				Refresh_Audio_Widgets ();
			}
		//Otherwise, just reset the necessary GUI components 	
		} else {
			Refresh_Audio_Widgets ();
			Stamp_GUI_Index = -1;
		}
		Refresh_Timestamp_Widgets ();
	}
	
	//Handler for clicking the timestamps
	public void TimeStampClick(object o,ButtonReleaseEventArgs args)
	{
		//Can't click if no corresponding audio track is selected
		if (Audio_GUI_Index < 1)
			return;
		//Change the current position of the Windows Player track if the corresponding mode is activated
		if (ActivationAsSearch) {
			MainClass.wmp.Ctlcontrols.currentPosition = lat [Audio_GUI_Index - 1].Timestamps [int.Parse (((Widget)o).Name.Split ('#') [1]) - 1].sec;
			return;
		}
		//Otherwise, reset index of the selected timestamp properly
		Stamp_GUI_Index = int.Parse (((Widget)o).Name.Split ('#') [1]) - 1;
		var times = lat [Audio_GUI_Index - 1].Timestamps;
		if (Stamp_GUI_Index > times.Count)
			Stamp_GUI_Index = -1;
		//Recover the text in the box based on the selected item value
		this.textbox.Text = (Stamp_GUI_Index > -1) ? times [Stamp_GUI_Index].title : "";
		//Reset GUI
		Refresh_Timestamp_Widgets ();
	}

	//Initialization of the main window
	public MainWindow (): base (Gtk.WindowType.Toplevel)
	{
		//Building GUI, launching timer to dynamically update the Windows Player state, setting base values
		Build ();
		GLib.Timeout.Add (250, new GLib.TimeoutHandler (PlayerTimer));
		textbox.ModifyFont (Pango.FontDescription.FromString("Monospace 20"));
		MainTitle.ModifyFont (Pango.FontDescription.FromString("Monospace 14"));
		textbox.Alignment = 0.5F;
		//Setting widget color scheme and splitter position
		SetWidgetColor (this, RGB2HEX (201, 201, 201));
		SetWidgetColor (eventbox20, RGB2HEX (235, 235, 235));
		SetWidgetColor (eventbox21, RGB2HEX (235, 235, 235));
		SetWidgetColor (eventbox22, RGB2HEX (235, 235, 235));
		SetWidgetColor (TitleBack, RGB2HEX (161, 161, 161));
		SetWidgetColor (PlayBack, RGB2HEX (77, 77, 77));
		SetWidgetColor (BlackBack, RGB2HEX (135, 135, 135));
		SetWidgetColor (WhiteBack, RGB2HEX (96, 96, 96));
		splitter1.Position = 400;
		//Morphing images into controls, attaching handlers for events
		ImagesAsControls ();
		ActivateImages ();
		//Loading dynamic widgets for correct visual representation of audiotracks and corresponding timestamps
		Dynamic_Widgets_Audio = FillLeftScrolledWindow (scrolledwindow1); Refresh_Audio_Widgets ();
		Dynamic_Widgets_Timestamps = FillRightScrolledWindow (scrolledwindow2);
		((Viewport)scrolledwindow1.Child).ShadowType=scrolledwindow1.ShadowType;
		((Viewport)scrolledwindow2.Child).ShadowType=scrolledwindow2.ShadowType;
	}

	//draw shadow at given coordinates x,y, (w)idth and (h)eight
	public void drawshadow(Cairo.Context context, int x, int y, int w, int h)
	{
		context.MoveTo (x + 1, y + h + 1);
		context.LineTo (w + 1, y + h + 1);
		context.LineTo (w + 1, y + 2);
		context.Stroke ();
	}

	//draw bar with specific (w)idth, (h)eight, grayscale (c)olor and fill level
	public void drawbar(Cairo.Context context, int w, int h, int c, double fill)
	{
		context.Save ();
		context.LineWidth = 0.75;
		//Fill the area
		double gray_value = ColorChan (c);
		context.SetSourceRGB (gray_value, gray_value, gray_value);
		context.Rectangle (0, 0, w, h);
		context.Fill ();
		//White color
		context.SetSourceRGB (1, 1, 1);
		//Updated coordinates with offset
		int h0 = 10, x0 = 10, w0 = w - x0, y0 = (h - 10) / 2, w1 = w0 - x0;
		context.Rectangle (x0, y0, w1, h0);
		context.Fill ();
		//Fill the bar
		double r = ColorChan (255), g = ColorChan (186), b = ColorChan (0);
		context.SetSourceRGB (r, g, b);
		context.Rectangle (x0, y0, w1 * fill, h0);
		context.Fill ();
		//Draw shadow of the bar
		gray_value = 0.2;
		for (int i = 0; i < 3; ++i) {
			context.SetSourceRGB (gray_value, gray_value, gray_value);
			drawshadow (context, x0 + i, y0 + i, w0 + i, h0 + i);
			context.LineWidth -= 0.2;
			gray_value += 0.3;
		}
		context.Restore ();
	}
	
	//2 handlers, drawing corresponding custom controls
	protected void OnSoundbarExposeEvent (object o, ExposeEventArgs args)
	{
		using (Cairo.Context context = Gdk.CairoHelper.Create (((Widget)o).GdkWindow)) {
			drawbar (context, args.Event.Area.Width, args.Event.Area.Height, 77, soundratio);
			context.GetTarget ().Dispose ();
			context.Dispose ();
		}
	}

	protected void OnWavebarExposeEvent (object o, ExposeEventArgs args)
	{
		using (Cairo.Context context = Gdk.CairoHelper.Create (((Widget)o).GdkWindow)) {
			drawbar (context, args.Event.Area.Width, args.Event.Area.Height, 235, playerratio);
			context.GetTarget ().Dispose ();
			context.Dispose ();
		}
	}
	
	//Functions soundmove and wavemove are used to set custom control position of audio sound and audio track, correspondingly,
	//based on provided mouse X-coordinate.
	void soundmove(int mouseX)
	{
		//Detect the measure of occupancy for the soundbar, as well as corresponding sound level
		int offset = 10;
		int zone = soundbar.Allocation.Width - offset * 2;
		int x_position = mouseX - offset;
		if (x_position < 0)
			x_position = 0;
		if (x_position > zone)
			x_position = zone;
		soundratio = x_position / (double)zone;
		//Redraw the soundbar
		soundbar.QueueDraw ();
		//Detect and set appropriate sound level (including muted sound scenario) 
		int soundlevel = ((int)(soundratio * 100.0));
		if (storedvol == -1)
			MainClass.wmp.settings.volume = soundlevel;
		else
			storedvol = soundlevel;
	}

	void wavemove(int mouseX)
	{
		//Same logic as in previous function
		int offset = 10;
		int zone = wavebar.Allocation.Width - offset * 2;
		int x_position = mouseX - offset;
		if (x_position < 0)
			x_position = 0;
		if (x_position > zone)
			x_position = zone;
		playerratio = x_position / (double)zone;
		wavebar.QueueDraw ();
		//If audiotrack is selected, reset the Windows Media Player current position
		if (Audio_GUI_Index > 0)
			MainClass.wmp.Ctlcontrols.currentPosition = MainClass.wmp.currentMedia.duration * playerratio;
	}		
	
	//Various handlers for custom controls (mouse motion & click)
	protected void OnSoundbarMotionNotifyEvent (object o, MotionNotifyEventArgs args)
	{
		soundmove ((int)args.Event.X);
	}

	protected void OnWavebarMotionNotifyEvent (object o, MotionNotifyEventArgs args)
	{
		wavemove ((int)args.Event.X);
	}	

	protected void OnSoundbarButtonPressEvent (object o, ButtonPressEventArgs args)
	{
		soundmove ((int)args.Event.X);
	}

	protected void OnWavebarButtonPressEvent (object o, ButtonPressEventArgs args)
	{
		wavemove ((int)args.Event.X);
	}
		
	//Handler for the main window closing (quit this application)
	protected void OnDeleteEvent (object sender, DeleteEventArgs a)
	{
		Application.Quit ();
		a.RetVal = true;
	}
}