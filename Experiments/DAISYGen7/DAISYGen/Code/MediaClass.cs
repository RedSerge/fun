using System;
using System.IO;
using System.Collections.Generic;
using Gtk;

namespace DAISYGen
{
	
	// class for storing timestamps:
	// seconds passed, title of timestamp, position in content hierarchy
	
	[Serializable]
	public class DaisyTime
	{
		public double sec;
		public string title;
		public int position;
		public DaisyTime(double s, string t, int p)
		{
			sec = s;
			title = t;
			position = p;
		}
	}

	// class for sorting timestamps
	
	public class TimeComparer:IComparer<DaisyTime>
	{
		public int Compare(DaisyTime first, DaisyTime second)
		{
			return Math.Sign (first.sec - second.sec);
		}
	}

	// class containing audiotracks with corresponding timestamps:
	// path to audiofile, length of audiotrack, timestamps

	[Serializable]
	public class AudioTracks
	{
		public string path;
		public double length = -1;
		public List<DaisyTime> Timestamps = new List<DaisyTime>();
		public AudioTracks(string path_to_audio)
		{
			path = path_to_audio;
			Timestamps.Add (new DaisyTime (0.0, "[Начало]", 0));
		}
	}

	// static class for displaying messages, Gtk# style
	 
	public static class Alert
	{
		public static void Show(Window win, string message, bool no_panic = false)
		{			
			var dialog = new MessageDialog (win, DialogFlags.Modal, no_panic?MessageType.Info:MessageType.Error, ButtonsType.Ok, message);
			dialog.Run ();
			dialog.Destroy ();
		}
	}
	
	// static class providing a method for exporting audiotracks in DAISY format;
	// pretty blunt technique of "paste-and-replace" kind, yet allowing to control the process more precisely.

	public static class ExportClass
	{
		//Template for sections
		const string hTemplate = "<h@ class=\"section\" id=\"rgn_cnt_*\"><a href=\"dtb_*.smil#rgn_txt_*_0001\">$</a></h@>";
		
		//Path to export
		public static string export_to = @"\Output";
		
		public static void Extract(Window win, List<AudioTracks> tracks)
		{
			//NCC file template
						
			List<string> ncc = new List<string> ();
			ncc.Add ("<?xml version=\"1.0\" encoding=\"utf-8\"?>");
			ncc.Add ("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/tr/xhtml1/dtd/xhtml1-transitional.dtd\">");
			ncc.Add ("<html xmlns=\"http://www.w3.org/1999/xhtml\">");
			ncc.Add ("<head>");
			ncc.Add ("<title>$</title>");
			ncc.Add ("<meta name=\"dc:format\" content=\"daisy 2.02\"/>");
			ncc.Add ("<meta http-equiv=\"content-type\" content=\"text/html; charset=utf-8\"/>");
			ncc.Add ("</head>");
			ncc.Add ("<body>");
			ncc.Add ("<h1 class=\"title\" id=\"rgn_cnt_0001\"><a href=\"dtb_0001.smil#rgn_txt_0001_0001\">$</a></h1>");
			ncc.Add ("</body>");
			ncc.Add ("</html>");
			
			//SMIL file template
			
			List<string> smil = new List<string> ();
			smil.Add ("<?xml version=\"1.0\" encoding=\"utf-8\"?>");
			smil.Add ("<!DOCTYPE smil PUBLIC \"-//W3C//DTD SMIL 1.0//EN\" \"http://www.w3.org/TR/REC-smil/SMIL10.dtd\">");
			smil.Add ("<smil>");
			smil.Add ("<head>");
			smil.Add ("<meta name=\"dc:format\" content=\"daisy 2.02\"/>");
			smil.Add ("<layout>");
			smil.Add ("<region id=\"txtview\"/>");
			smil.Add ("</layout>");
			smil.Add ("</head>");
			smil.Add ("<body>");
			smil.Add ("<seq>");
			smil.Add ("<par endsync=\"last\">");
			smil.Add ("<text src=\"ncc.html#rgn_cnt_*\" id=\"rgn_txt_*_0001\"/>");
			smil.Add ("<seq>");
			smil.Add ("<audio src=\"$\" clip-begin=\"npt=@s\" clip-end=\"npt=%s\" id=\"phr_*_0001\"/>");
			smil.Add ("</seq>");
			smil.Add ("</par>");
			smil.Add ("</seq>");
			smil.Add ("</body>");
			smil.Add ("</smil>");
			
			//Merging the template into one string
			string smil_join = String.Join (Environment.NewLine, smil);
			
			//Initialization
			int sum = 0;
			string export_path = Path.GetDirectoryName (MainClass.path) + export_to;
			bool firstpass = true;
			string title = string.Empty;
			foreach (AudioTracks t in tracks)
			{
				//Check if track is actually valid
				if (t.length <= 0)
				{
					Alert.Show (win, "Track " + t.path + " has not been initialized. Activate it again or make sure its length > 0 ms");
					return;
				}
				
				//sum represents total count of timestamps overall
				sum += t.Timestamps.Count;
				
				//The very first timestamp title is actually added to the corresponding NCC tag later 
				if (sum!=0 && firstpass)
				{
					title = t.Timestamps[0].title;
					firstpass = false;
				}
			}
			
			//Check if book is not empty
			if (sum == 0)
			{
				Alert.Show (win, "Can't export empty book");
				return;
			}
			
			//Clean the output directory (or make it) and correct the path
			if (Directory.Exists (export_path))
				foreach (string filename in Directory.GetFiles(export_path))
					File.Delete (filename);
			else
				Directory.CreateDirectory (export_path);
			export_path += @"\";
			
			//Buffer to write for further export of NCC primary file
			var write = new List<string> ();
			write.AddRange(ncc);
			
			//Pre-format; title inclusion, audiotrack identification unification
			write[4]=write[4].Replace("$", title);
			write[9]=write[9].Replace("$", title);
			for(int i = sum; i > 1; --i)
			{
				write.Insert(10, hTemplate.Replace("*", i.ToString("0000")));
			}
			
			//The most important part...
			int index = 9, track_count = tracks.Count;
			for(int i=0; i < track_count; ++i)
			{
				//Audiotrack depersonalization for further unification
				string extended_track_path=(i+1).ToString()+Path.GetExtension(tracks[i].path);
				File.Copy(tracks[i].path, export_path + extended_track_path);
				
				//Parsing corresponding timestamps
				int stamp_count = tracks[i].Timestamps.Count;
				for(int j = 0; j < stamp_count; ++j)
				{
					DaisyTime CurrentTime=tracks[i].Timestamps[j];
					
					//Paste information about timestamp, if necessary
					if (index != 9)
						write[index]=write[index].Replace("$", CurrentTime.title).Replace("@", (CurrentTime.position+1).ToString());
					
					//SMIL file generation for every timestamp
					double subexpression = (j == stamp_count-1) ? tracks[i].length : tracks[i].Timestamps[j + 1].sec;
					string subexpression_for_file = (index - 8).ToString("0000");
					File.WriteAllText(export_path + "dtb_" + subexpression_for_file + ".smil", smil_join.Replace("@", CurrentTime.sec.ToString("0.000")).Replace("%", subexpression.ToString("0.000")).Replace("*", subexpression_for_file).Replace("$", extended_track_path));
					++index;
				}
			}
			
			//Export NCC file and confirm the export finalization
			File.WriteAllLines (export_path + "ncc.html", write);
			Alert.Show (win, "Done", true);
		}
	}
}